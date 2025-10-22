package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateRecurringTransaction(recurringTransaction *models.RecurringTransaction, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO recurring_transactions (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", models.RecurringTransactionColumns)
	_, err := db.Exec(query, recurringTransaction.ID, recurringTransaction.UserID, recurringTransaction.AccountID, recurringTransaction.CategoryID, recurringTransaction.Description, recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.RecurringFrequency, recurringTransaction.RecurringDate, recurringTransaction.CreatedAt, recurringTransaction.UpdatedAt)
	return err
}

func GetRecurringTransactionsByUserID(userID uuid.UUID, db *sql.DB) ([]models.RecurringTransaction, error) {
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurringTransactions []models.RecurringTransaction
	for rows.Next() {
		var recurringTransaction models.RecurringTransaction
		if err := rows.Scan(&recurringTransaction.ID, &recurringTransaction.UserID, &recurringTransaction.AccountID, &recurringTransaction.CategoryID, &recurringTransaction.Description, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.RecurringFrequency, &recurringTransaction.RecurringDate, &recurringTransaction.CreatedAt, &recurringTransaction.UpdatedAt); err != nil {
			return nil, err
		}
		recurringTransactions = append(recurringTransactions, recurringTransaction)
	}
	return recurringTransactions, nil
}

func GetRecurringTransactionByID(id uuid.UUID, db *sql.DB) (*models.RecurringTransaction, error) {
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions WHERE id = $1"
	row := db.QueryRow(query, id)

	var recurringTransaction models.RecurringTransaction
	if err := row.Scan(&recurringTransaction.ID, &recurringTransaction.UserID, &recurringTransaction.AccountID, &recurringTransaction.CategoryID, &recurringTransaction.Description, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.RecurringFrequency, &recurringTransaction.RecurringDate, &recurringTransaction.CreatedAt, &recurringTransaction.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &recurringTransaction, nil
}

func UpdateRecurringTransaction(recurringTransaction *models.RecurringTransaction, db *sql.DB) error {
	query := "UPDATE recurring_transactions SET account_id = $1, category_id = $2, description = $3, amount = $4, type = $5, recurring_frequency = $6, recurring_date = $7, updated_at = $8 WHERE id = $9"
	_, err := db.Exec(query, recurringTransaction.AccountID, recurringTransaction.CategoryID, recurringTransaction.Description, recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.RecurringFrequency, recurringTransaction.RecurringDate, recurringTransaction.UpdatedAt, recurringTransaction.ID)
	return err
}

func DeleteRecurringTransaction(id uuid.UUID, db *sql.DB) error {
	query := "DELETE FROM recurring_transactions WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
