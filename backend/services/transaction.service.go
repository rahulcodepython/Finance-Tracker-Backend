package services

import (
	"database/sql"
	"errors"
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
		return nil, errors.New("Category not found")
	}

	transactionType := models.TransactionType(category.Type)

	// Update account balance
	account, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("Account not found")
	}

	if transactionType == models.TransactionTypeIncome {
		account.Balance += amount
	} else {
		account.Balance -= amount
	}
	if err := repository.UpdateAccount(account, db); err != nil {
		return nil, err
	}

	// Update budget if provided
	if budgetID.Valid {
		budget, err := repository.GetBudgetByID(budgetID.UUID, db)
		if err != nil {
			return nil, err
		}

		if budget == nil {
			return nil, errors.New("Budget not found")
		}

		budget.Amount -= amount // Deduct new transaction amount from new budget
		if err := repository.UpdateBudget(budget, db); err != nil {
			return nil, err
		}
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

	// Create the transaction
	if err := repository.CreateTransaction(transaction, db); err != nil {
		return nil, err
	}

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
		return nil, errors.New("Transaction not found")
	}

	if transaction.CategoryID != categoryID {
		category, err := repository.GetCategoryByID(categoryID, db)
		if err != nil {
			return nil, err
		}

		if category == nil {
			return nil, errors.New("Category not found")
		}
	}

	if transaction.AccountID != accountID {
		oldAccountID := transaction.AccountID
		oldAmount := transaction.Amount
		oldType := transaction.Type

		// Revert old transaction amount from old account balance
		oldAccount, err := repository.GetAccountByID(oldAccountID, db)
		if err != nil {
			return nil, err
		}
		if oldType == models.TransactionTypeIncome {
			oldAccount.Balance -= oldAmount
		} else {
			oldAccount.Balance += oldAmount
		}
		if err := repository.UpdateAccount(oldAccount, db); err != nil {
			return nil, err
		}

		// Update new account balance
		newAccount, err := repository.GetAccountByID(accountID, db)
		if err != nil {
			return nil, err
		}

		if newAccount == nil {
			return nil, errors.New("Account not found")
		}

		if transaction.Type == models.TransactionTypeIncome {
			newAccount.Balance += amount
		} else {
			newAccount.Balance -= amount
		}
		if err := repository.UpdateAccount(newAccount, db); err != nil {
			return nil, err
		}
	} else if transaction.Amount != amount {
		difference := amount - transaction.Amount

		account, err := repository.GetAccountByID(transaction.AccountID, db)
		if err != nil {
			return nil, err
		}
		if transaction.Type == models.TransactionTypeIncome {
			account.Balance += difference
		} else {
			account.Balance -= difference
		}
		if err := repository.UpdateAccount(account, db); err != nil {
			return nil, err
		}
	}

	if transaction.BudgetID.Valid && transaction.BudgetID.UUID != budgetID.UUID {
		oldBudgetID := transaction.BudgetID
		oldAmount := transaction.Amount

		// Revert old transaction amount from old budget if it existed
		if oldBudgetID.Valid {
			oldBudget, err := repository.GetBudgetByID(oldBudgetID.UUID, db)
			if err != nil {
				return nil, err
			}
			oldBudget.Amount += oldAmount // Add back the old amount
			if err := repository.UpdateBudget(oldBudget, db); err != nil {
				return nil, err
			}
		}

		// Update new budget if provided
		if budgetID.Valid {
			newBudget, err := repository.GetBudgetByID(budgetID.UUID, db)
			if err != nil {
				return nil, err
			}
			if newBudget == nil {
				return nil, errors.New("Budget not found")
			}
			newBudget.Amount -= amount // Deduct new transaction amount from new budget
			if err := repository.UpdateBudget(newBudget, db); err != nil {
				return nil, err
			}
		}
	} else if transaction.Amount != amount {
		difference := amount - transaction.Amount

		if budgetID.Valid {
			budget, err := repository.GetBudgetByID(transaction.BudgetID.UUID, db)
			if err != nil {
				return nil, err
			}
			budget.Amount -= difference // Deduct new transaction amount from new budget
			if err := repository.UpdateBudget(budget, db); err != nil {
				return nil, err
			}
		}
	}

	// Update transaction details
	transaction.AccountID = accountID
	transaction.CategoryID = categoryID
	transaction.BudgetID = budgetID
	transaction.Description = description
	transaction.Amount = amount
	transaction.TransactionDate = transactionDate
	transaction.Note = note
	transaction.UpdatedAt = time.Now().In(utils.LOC)

	if err := repository.UpdateTransaction(transaction, db); err != nil {
		return nil, err
	}

	return transaction, nil
}

func DeleteTransaction(id uuid.UUID, db *sql.DB) error {
	transaction, err := repository.GetTransactionByID(id, db)
	if err != nil {
		return err
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

	if err := repository.UpdateAccount(account, db); err != nil {
		return err
	}

	// Revert budget if it existed
	if transaction.BudgetID.Valid {
		budget, err := repository.GetBudgetByID(transaction.BudgetID.UUID, db)
		if err != nil {
			return err
		}
		budget.Amount += transaction.Amount // Add back the transaction amount
		if err := repository.UpdateBudget(budget, db); err != nil {
			return err
		}
	}

	if err := repository.DeleteTransaction(id, db); err != nil {
		return err
	}

	return nil
}

func GetAggregateData(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	return repository.GetAggregateDataByUserID(userID, startDate, endDate, db)
}

func GetSpendingByCategory(userID uuid.UUID, db *sql.DB) ([]map[string]interface{}, error) {
	return repository.GetSpendingByCategory(userID, db)
}
