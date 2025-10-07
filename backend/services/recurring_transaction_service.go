package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateRecurringTransaction(userID uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, description string, amount float64, transactionType models.TransactionType, frequency string, startDate time.Time, endDate time.Time, db *sql.DB) (*models.RecurringTransaction, error) {
	return nil, nil
}

func GetRecurringTransactions(userID uuid.UUID, db *sql.DB) ([]models.RecurringTransaction, error) {
	return nil, nil
}

func UpdateRecurringTransaction(id uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, description string, amount float64, transactionType models.TransactionType, frequency string, startDate time.Time, endDate time.Time, db *sql.DB) (*models.RecurringTransaction, error) {
	return nil, nil
}

func DeleteRecurringTransaction(id uuid.UUID, db *sql.DB) error {
	return nil
}
