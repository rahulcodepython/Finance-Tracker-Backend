package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func CreateRecurringTransaction(userID uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, description string, amount float64, transactionType models.TransactionType, recurringFrequency models.RecurringFrequency, recurringDate int, db *sql.DB) (*models.RecurringTransaction, error) {
	recurringTransaction := &models.RecurringTransaction{
		ID:                 uuid.New(),
		UserID:             userID,
		AccountID:          accountID,
		CategoryID:         uuid.NullUUID{UUID: categoryID, Valid: true},
		Description:        description,
		Amount:             amount,
		Type:               transactionType,
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

func UpdateRecurringTransaction(id uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, description string, amount float64, transactionType models.TransactionType, recurringFrequency models.RecurringFrequency, recurringDate int, db *sql.DB) (*models.RecurringTransaction, error) {
	recurringTransaction, err := repository.GetRecurringTransactionByID(id, db)
	if err != nil {
		return nil, err
	}

	recurringTransaction.AccountID = accountID
	recurringTransaction.CategoryID = uuid.NullUUID{UUID: categoryID, Valid: true}
	recurringTransaction.Description = description
	recurringTransaction.Amount = amount
	recurringTransaction.Type = transactionType
	recurringTransaction.RecurringFrequency = recurringFrequency
	recurringTransaction.RecurringDate = recurringDate
	recurringTransaction.UpdatedAt = time.Now().In(utils.LOC)

	if err := repository.UpdateRecurringTransaction(recurringTransaction, db); err != nil {
		return nil, err
	}

	return recurringTransaction, nil
}

func DeleteRecurringTransaction(id uuid.UUID, db *sql.DB) error {
	return repository.DeleteRecurringTransaction(id, db)
}
