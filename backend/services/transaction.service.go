package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func CreateTransaction(userID uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, transactionDate time.Time, note sql.NullString, db *sql.DB) (*models.Transaction, error) {
	category, err := repository.GetCategoryByID(categoryID, db)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, sql.ErrNoRows
	}

	transactionType := models.TransactionType(category.Type)

	// Update account balance
	account, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, sql.ErrNoRows
	}

	if transactionType == models.TransactionTypeIncome {
		account.Balance += amount
	} else {
		account.Balance -= amount
	}

	// Update budget if provided
	var budget *models.Budget

	if budgetID.Valid {
		budget, err = repository.GetBudgetByID(budgetID.UUID, db)
		if err != nil {
			return nil, err
		}

		if budget == nil {
			return nil, sql.ErrNoRows
		}

		budget.Amount -= amount // Deduct new transaction amount from new budget
	}

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

	err = utils.DBTransaction(db, func(tx *sql.Tx) error {
		if err := repository.UpdateAccount(account, tx); err != nil {
			return err
		}

		if err := repository.UpdateBudget(budget, tx); err != nil {
			return err
		}

		// Create the transaction
		if err := repository.CreateTransaction(transaction, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Log the creation
	go CreateLog(userID, fmt.Sprintf("New transaction '%s' created", transaction.Description), db)

	return transaction, nil
}

func GetTransactions(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, budgetID string, startDate string, endDate string, db *sql.DB) ([]models.Transaction, error) {
	return repository.GetTransactionsByUserIDWithFilters(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
}

func UpdateTransaction(id uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, transactionDate time.Time, note sql.NullString, db *sql.DB) (*models.Transaction, error) {
	transaction, err := repository.GetTransactionByID(id, db)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, sql.ErrNoRows
	}

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

	if transaction.AccountID != accountID {
		oldAccountID := transaction.AccountID
		oldAmount := transaction.Amount
		oldType := transaction.Type

		// Get old account
		oldAccount, err := repository.GetAccountByID(oldAccountID, db)
		if err != nil {
			return nil, err
		}
		if oldAccount == nil {
			return nil, sql.ErrNoRows
		}
		// revert old transaction amount from old account balance
		if oldType == models.TransactionTypeIncome {
			oldAccount.Balance -= oldAmount
		} else {
			oldAccount.Balance += oldAmount
		}
		oldAccount.UpdatedAt = time.Now().In(utils.LOC)
		oldAccountToUpdate = oldAccount

		// Get new account
		newAccount, err := repository.GetAccountByID(accountID, db)
		if err != nil {
			return nil, err
		}
		if newAccount == nil {
			return nil, sql.ErrNoRows
		}
		// apply new amount to new account according to transaction.Type
		if transaction.Type == models.TransactionTypeIncome {
			newAccount.Balance += amount
		} else {
			newAccount.Balance -= amount
		}
		newAccount.UpdatedAt = time.Now().In(utils.LOC)
		newAccountToUpdate = newAccount
	} else if transaction.Amount != amount {
		// Same account; only adjust by difference
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

	if transaction.BudgetID.Valid && (!budgetID.Valid || transaction.BudgetID.UUID != budgetID.UUID) {
		oldBudgetID := transaction.BudgetID
		oldAmount := transaction.Amount

		// Revert old budget amount
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

	transaction.AccountID = accountID
	transaction.CategoryID = categoryID
	transaction.BudgetID = budgetID
	transaction.Description = description
	transaction.Amount = amount
	transaction.TransactionDate = transactionDate
	transaction.Note = note
	transaction.UpdatedAt = time.Now().In(utils.LOC)

	err = utils.DBTransaction(db, func(tx *sql.Tx) error {
		if oldAccountToUpdate != nil {
			if err := repository.UpdateAccount(oldAccountToUpdate, tx); err != nil {
				return err
			}
		}
		if newAccountToUpdate != nil {
			if err := repository.UpdateAccount(newAccountToUpdate, tx); err != nil {
				return err
			}
		}
		if accountToUpdate != nil {
			if err := repository.UpdateAccount(accountToUpdate, tx); err != nil {
				return err
			}
		}

		if oldBudgetToUpdate != nil {
			if err := repository.UpdateBudget(oldBudgetToUpdate, tx); err != nil {
				return err
			}
		}
		if newBudgetToUpdate != nil {
			if err := repository.UpdateBudget(newBudgetToUpdate, tx); err != nil {
				return err
			}
		}

		if err := repository.UpdateTransaction(transaction, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Log the update
	go CreateLog(transaction.UserID, fmt.Sprintf("Transaction '%s' updated", transaction.Description), db)

	return transaction, nil

}

func DeleteTransaction(id uuid.UUID, db *sql.DB) error {
	transaction, err := repository.GetTransactionByID(id, db)
	if err != nil {
		return err
	}

	if transaction == nil {
		return sql.ErrNoRows
	}

	account, err := repository.GetAccountByID(transaction.AccountID, db)
	if err != nil {
		return err
	}

	if transaction.Type == models.TransactionTypeIncome {
		account.Balance -= transaction.Amount
	} else {
		account.Balance += transaction.Amount
	}
	account.UpdatedAt = time.Now().In(utils.LOC)

	// Revert budget if it existed
	var budget *models.Budget

	if transaction.BudgetID.Valid {
		budget, err = repository.GetBudgetByID(transaction.BudgetID.UUID, db)
		if err != nil {
			return err
		}
		budget.Amount += transaction.Amount // Add back the transaction amount
		budget.UpdatedAt = time.Now().In(utils.LOC)
	}

	err = utils.DBTransaction(db, func(tx *sql.Tx) error {
		if err := repository.UpdateAccount(account, tx); err != nil {
			return err
		}

		if err := repository.UpdateBudget(budget, tx); err != nil {
			return err
		}

		// Create the transaction
		if err := repository.DeleteTransaction(id, db); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Log the deletion
	go CreateLog(transaction.UserID, fmt.Sprintf("Transaction '%s' removed", transaction.Description), db)

	return nil
}

func GetAggregateData(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	return repository.GetAggregateDataByUserID(userID, startDate, endDate, db)
}

func GetSpendingByCategory(userID uuid.UUID, db *sql.DB) ([]map[string]interface{}, error) {
	return repository.GetSpendingByCategory(userID, db)
}
