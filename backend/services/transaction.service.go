// Package services provides business logic for the application.
package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/serializers"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateTransaction creates a new transaction for a user.
func CreateTransaction(userID uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, transactionDate time.Time, note sql.NullString, db *sql.DB) (*serializers.TransactionResponse, error) {
	// Get the category to determine the transaction type.
	category, err := repository.GetCategoryByID(categoryID, db)
	if err != nil {
		return nil, err
	}

	// If the category does not exist, return an error.
	if category == nil {
		return nil, sql.ErrNoRows
	}

	// Set the transaction type based on the category type.
	transactionType := models.TransactionType(category.Type)

	// Get the account to update its balance.
	account, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}

	// If the account does not exist, return an error.
	if account == nil {
		return nil, sql.ErrNoRows
	}

	// Update the account balance based on the transaction type.
	if transactionType == models.TransactionTypeIncome {
		account.Balance += amount
	} else {
		account.Balance -= amount
	}

	// If a budget is provided, update its amount.
	var budget *models.Budget
	if budgetID.Valid {
		budget, err = repository.GetBudgetByID(budgetID.UUID, db)
		if err != nil {
			return nil, err
		}

		if budget == nil {
			return nil, sql.ErrNoRows
		}

		budget.Amount -= amount
	}

	// Create a new Transaction model.
	transaction := &models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		AccountID:       accountID,
		CategoryID:      categoryID,
		BudgetID:        budgetID,
		Description:     description,
		Amount:          amount,
		Type:            transactionType,
		TransactionDate: transactionDate,
		Note:            note,
		CreatedAt:       time.Now().In(utils.LOC),
		UpdatedAt:       time.Now().In(utils.LOC),
	}

	// Use a database transaction to ensure atomicity.
	err = utils.DBTransaction(db, func(tx *sql.Tx) error {
		// Update the account balance.
		if err := repository.UpdateAccount(account, tx); err != nil {
			return err
		}

		// If a budget was updated, save the changes.
		if budget != nil {
			if err := repository.UpdateBudget(budget, tx); err != nil {
				return err
			}
		}

		// Create the transaction.
		if err := repository.CreateTransaction(transaction, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Create a log entry for the transaction creation.
	go CreateLog(userID, fmt.Sprintf("New transaction '%s' created", transaction.Description), db)

	var transactionResponse serializers.TransactionResponse
	transactionResponse.ID = transaction.ID
	transactionResponse.UserID = transaction.UserID
	transactionResponse.Description = transaction.Description
	transactionResponse.Amount = transaction.Amount
	transactionResponse.Type = transaction.Type
	transactionResponse.TransactionDate = transaction.TransactionDate
	transactionResponse.Note = transaction.Note
	transactionResponse.CreatedAt = transaction.CreatedAt
	transactionResponse.UpdatedAt = transaction.UpdatedAt
	transactionResponse.Account.UUID = transaction.AccountID
	transactionResponse.Category.UUID = transaction.CategoryID
	if transaction.BudgetID.Valid {
		transactionResponse.Budget.UUID = transaction.BudgetID.UUID
	}

	return &transactionResponse, nil
}

// GetTransactions retrieves all transactions for a user with optional filters and pagination.
func GetTransactions(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, budgetID string, startDate string, endDate string, db *sql.DB) ([]serializers.TransactionResponse, error) {
	return repository.GetTransactionsByUserIDWithFilters(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
}

// UpdateTransaction updates an existing transaction.
func UpdateTransaction(id uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, transactionDate time.Time, note sql.NullString, db *sql.DB) (*serializers.TransactionResponse, error) {
	// Get the original transaction from the database.
	transaction, err := repository.GetTransactionByID(id, db)
	if err != nil {
		return nil, err
	}

	// If the transaction does not exist, return an error.
	if transaction == nil {
		return nil, sql.ErrNoRows
	}

	// If the category has changed, get the new category.
	if transaction.CategoryID != categoryID {
		category, err := repository.GetCategoryByID(categoryID, db)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, sql.ErrNoRows
		}
	}

	var (
		oldAccountToUpdate *models.Account
		newAccountToUpdate *models.Account
		accountToUpdate    *models.Account
		oldBudgetToUpdate  *models.Budget
		newBudgetToUpdate  *models.Budget
	)

	// If the account has changed, revert the old account balance and update the new one.
	if transaction.AccountID != accountID {
		oldAccountID := transaction.AccountID
		oldAmount := transaction.Amount
		oldType := transaction.Type

		// Get the old account.
		oldAccount, err := repository.GetAccountByID(oldAccountID, db)
		if err != nil {
			return nil, err
		}
		if oldAccount == nil {
			return nil, sql.ErrNoRows
		}
		// Revert the old transaction amount from the old account balance.
		if oldType == models.TransactionTypeIncome {
			oldAccount.Balance -= oldAmount
		} else {
			oldAccount.Balance += oldAmount
		}
		oldAccount.UpdatedAt = time.Now().In(utils.LOC)
		oldAccountToUpdate = oldAccount

		// Get the new account.
		newAccount, err := repository.GetAccountByID(accountID, db)
		if err != nil {
			return nil, err
		}
		if newAccount == nil {
			return nil, sql.ErrNoRows
		}
		// Apply the new amount to the new account according to the transaction type.
		if transaction.Type == models.TransactionTypeIncome {
			newAccount.Balance += amount
		} else {
			newAccount.Balance -= amount
		}
		newAccount.UpdatedAt = time.Now().In(utils.LOC)
		newAccountToUpdate = newAccount
	} else if transaction.Amount != amount {
		// If the account is the same but the amount has changed, adjust the balance by the difference.
		difference := amount - transaction.Amount
		account, err := repository.GetAccountByID(transaction.AccountID, db)
		if err != nil {
			return nil, err
		}
		if account == nil {
			return nil, sql.ErrNoRows
		}
		if transaction.Type == models.TransactionTypeIncome {
			account.Balance += difference
		} else {
			account.Balance -= difference
		}
		account.UpdatedAt = time.Now().In(utils.LOC)
		accountToUpdate = account
	}

	// If the budget has changed, revert the old budget amount and update the new one.
	if transaction.BudgetID.Valid && (!budgetID.Valid || transaction.BudgetID.UUID != budgetID.UUID) {
		oldBudgetID := transaction.BudgetID
		oldAmount := transaction.Amount

		// Revert the old budget amount.
		if oldBudgetID.Valid {
			oldBudget, err := repository.GetBudgetByID(oldBudgetID.UUID, db)
			if err != nil {
				return nil, err
			}
			if oldBudget == nil {
				return nil, sql.ErrNoRows
			}
			oldBudget.Amount += oldAmount
			oldBudget.UpdatedAt = time.Now().In(utils.LOC)
			oldBudgetToUpdate = oldBudget
		}

		// If a new budget is provided, update its amount.
		if budgetID.Valid {
			newBudget, err := repository.GetBudgetByID(budgetID.UUID, db)
			if err != nil {
				return nil, err
			}
			if newBudget == nil {
				return nil, sql.ErrNoRows
			}
			newBudget.Amount -= amount
			newBudget.UpdatedAt = time.Now().In(utils.LOC)
			newBudgetToUpdate = newBudget
		}
	} else if transaction.Amount != amount {
		// If the budget is the same but the amount has changed, adjust the budget amount by the difference.
		if budgetID.Valid {
			budget, err := repository.GetBudgetByID(transaction.BudgetID.UUID, db)
			if err != nil {
				return nil, err
			}
			if budget == nil {
				return nil, sql.ErrNoRows
			}
			difference := amount - transaction.Amount
			budget.Amount -= difference
			budget.UpdatedAt = time.Now().In(utils.LOC)
			oldBudgetToUpdate = budget
		}
	}

	// Update the transaction fields.
	transaction.AccountID = accountID
	transaction.CategoryID = categoryID
	transaction.BudgetID = budgetID
	transaction.Description = description
	transaction.Amount = amount
	transaction.TransactionDate = transactionDate
	transaction.Note = note
	transaction.UpdatedAt = time.Now().In(utils.LOC)

	// Use a database transaction to ensure atomicity.
	err = utils.DBTransaction(db, func(tx *sql.Tx) error {
		// Update the old account if it exists.
		if oldAccountToUpdate != nil {
			if err := repository.UpdateAccount(oldAccountToUpdate, tx); err != nil {
				return err
			}
		}
		// Update the new account if it exists.
		if newAccountToUpdate != nil {
			if err := repository.UpdateAccount(newAccountToUpdate, tx); err != nil {
				return err
			}
		}
		// Update the account if it exists.
		if accountToUpdate != nil {
			if err := repository.UpdateAccount(accountToUpdate, tx); err != nil {
				return err
			}
		}

		// Update the old budget if it exists.
		if oldBudgetToUpdate != nil {
			if err := repository.UpdateBudget(oldBudgetToUpdate, tx); err != nil {
				return err
			}
		}
		// Update the new budget if it exists.
		if newBudgetToUpdate != nil {
			if err := repository.UpdateBudget(newBudgetToUpdate, tx); err != nil {
				return err
			}
		}

		// Update the transaction.
		if err := repository.UpdateTransaction(transaction, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Create a log entry for the transaction update.
	go CreateLog(transaction.UserID, fmt.Sprintf("Transaction '%s' updated", transaction.Description), db)

	var transactionResponse serializers.TransactionResponse
	transactionResponse.ID = transaction.ID
	transactionResponse.UserID = transaction.UserID
	transactionResponse.Description = transaction.Description
	transactionResponse.Amount = transaction.Amount
	transactionResponse.Type = transaction.Type
	transactionResponse.TransactionDate = transaction.TransactionDate
	transactionResponse.Note = transaction.Note
	transactionResponse.CreatedAt = transaction.CreatedAt
	transactionResponse.UpdatedAt = transaction.UpdatedAt
	transactionResponse.Account.UUID = transaction.AccountID
	transactionResponse.Category.UUID = transaction.CategoryID
	if transaction.BudgetID.Valid {
		transactionResponse.Budget.UUID = transaction.BudgetID.UUID
	}

	return &transactionResponse, nil
}

// DeleteTransaction deletes a transaction.
func DeleteTransaction(id uuid.UUID, db *sql.DB) error {
	// Get the transaction from the database.
	transaction, err := repository.GetTransactionByID(id, db)
	if err != nil {
		return err
	}

	// If the transaction does not exist, return an error.
	if transaction == nil {
		return sql.ErrNoRows
	}

	// Get the account to revert its balance.
	account, err := repository.GetAccountByID(transaction.AccountID, db)
	if err != nil {
		return err
	}

	// Revert the account balance based on the transaction type.
	if transaction.Type == models.TransactionTypeIncome {
		account.Balance -= transaction.Amount
	} else {
		account.Balance += transaction.Amount
	}
	account.UpdatedAt = time.Now().In(utils.LOC)

	// If a budget was associated with the transaction, revert its amount.
	var budget *models.Budget
	if transaction.BudgetID.Valid {
		budget, err = repository.GetBudgetByID(transaction.BudgetID.UUID, db)
		if err != nil {
			return err
		}
		budget.Amount += transaction.Amount
		budget.UpdatedAt = time.Now().In(utils.LOC)
	}

	// Use a database transaction to ensure atomicity.
	err = utils.DBTransaction(db, func(tx *sql.Tx) error {
		// Update the account balance.
		if err := repository.UpdateAccount(account, tx); err != nil {
			return err
		}

		// If a budget was updated, save the changes.
		if budget != nil {
			if err := repository.UpdateBudget(budget, tx); err != nil {
				return err
			}
		}

		// Delete the transaction.
		if err := repository.DeleteTransaction(id, db); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Create a log entry for the transaction deletion.
	go CreateLog(transaction.UserID, fmt.Sprintf("Transaction '%s' removed", transaction.Description), db)

	return nil
}

// GetAggregateData retrieves aggregate transaction data for a user.
func GetAggregateData(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (*serializers.DashboardSummary, error) {
	return repository.GetAggregateDataByUserID(userID, startDate, endDate, db)
}

// GetSpendingByCategory retrieves the total spending by category for a user.
func GetSpendingByCategory(userID uuid.UUID, db *sql.DB) (*[]serializers.CategoryAggregate, error) {
	return repository.GetSpendingByCategory(userID, db)
}

func GetIncomeExpense(userID uuid.UUID, db *sql.DB) (*serializers.IncomeExpenseAggregate, error) {
	return repository.GetIncomeExpense(userID, db)
}
