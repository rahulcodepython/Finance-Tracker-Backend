package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func CreateRecurringTransaction(userID uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, note sql.NullString, recurringFrequency models.RecurringFrequency, recurringDate int, db *sql.DB) (*models.RecurringTransaction, error) {
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

	// Update budget if provided
	if budgetID.Valid {
		budget, err := repository.GetBudgetByID(budgetID.UUID, db)
		if err != nil {
			return nil, err
		}

		if budget == nil {
			return nil, sql.ErrNoRows
		}
	}

	recurringTransaction := &models.RecurringTransaction{
		ID:                 uuid.New(),
		UserID:             userID,
		AccountID:          accountID,
		CategoryID:         categoryID,
		BudgetID:           budgetID,
		Description:        description,
		Amount:             amount,
		Type:               transactionType,
		Note:               note,
		RecurringFrequency: recurringFrequency,
		RecurringDate:      recurringDate,
		CreatedAt:          time.Now().In(utils.LOC),
		UpdatedAt:          time.Now().In(utils.LOC),
	}

	if err := repository.CreateRecurringTransaction(recurringTransaction, db); err != nil {
		return nil, err
	}

	return recurringTransaction, nil
}

func GetRecurringTransactions(userID uuid.UUID, db *sql.DB) ([]models.RecurringTransaction, error) {
	return repository.GetRecurringTransactionsByUserID(userID, db)
}

func UpdateRecurringTransaction(id uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, note sql.NullString, recurringFrequency models.RecurringFrequency, recurringDate int, db *sql.DB) (*models.RecurringTransaction, error) {
	recurringTransaction, err := repository.GetRecurringTransactionByID(id, db)
	if err != nil {
		return nil, err
	}

	if recurringTransaction == nil {
		return nil, sql.ErrNoRows
	}

	category, err := repository.GetCategoryByID(categoryID, db)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, sql.ErrNoRows
	}

	transactionType := models.TransactionType(category.Type)

	account, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, sql.ErrNoRows
	}

	if budgetID.Valid {
		budget, err := repository.GetBudgetByID(budgetID.UUID, db)
		if err != nil {
			return nil, err
		}

		if budget == nil {
			return nil, sql.ErrNoRows
		}
	}

	recurringTransaction.AccountID = accountID
	recurringTransaction.CategoryID = categoryID
	recurringTransaction.BudgetID = budgetID
	recurringTransaction.Description = description
	recurringTransaction.Amount = amount
	recurringTransaction.Type = transactionType
	recurringTransaction.Note = note
	recurringTransaction.RecurringFrequency = recurringFrequency
	recurringTransaction.RecurringDate = recurringDate
	recurringTransaction.UpdatedAt = time.Now().In(utils.LOC)

	if err := repository.UpdateRecurringTransaction(recurringTransaction, db); err != nil {
		return nil, err
	}

	return recurringTransaction, nil
}

func DeleteRecurringTransaction(id uuid.UUID, db *sql.DB) error {
	recurringTransaction, err := repository.GetRecurringTransactionByID(id, db)
	if err != nil {
		return err
	}

	if recurringTransaction == nil {
		return sql.ErrNoRows
	}

	return repository.DeleteRecurringTransaction(id, db)
}
