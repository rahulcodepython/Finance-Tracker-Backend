package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateRecurringTransaction(recurringTransaction *models.RecurringTransaction, db interfaces.SqlExecutor) error {
	query := fmt.Sprintf("INSERT INTO recurring_transactions (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", models.RecurringTransactionColumns)
	_, err := db.Exec(query, recurringTransaction.ID, recurringTransaction.UserID, recurringTransaction.AccountID, recurringTransaction.CategoryID, recurringTransaction.BudgetID, recurringTransaction.Description, recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.Note, recurringTransaction.RecurringFrequency, recurringTransaction.RecurringDate, recurringTransaction.CreatedAt, recurringTransaction.UpdatedAt)
	return err
}

func GetRecurringTransactionsByUserID(userID uuid.UUID, db interfaces.SqlExecutor) ([]models.RecurringTransaction, error) {
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurringTransactions []models.RecurringTransaction
	for rows.Next() {
		var recurringTransaction models.RecurringTransaction
		if err := rows.Scan(&recurringTransaction.ID, &recurringTransaction.UserID, &recurringTransaction.AccountID, &recurringTransaction.CategoryID, &recurringTransaction.BudgetID, &recurringTransaction.Description, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.Note, &recurringTransaction.RecurringFrequency, &recurringTransaction.RecurringDate, &recurringTransaction.CreatedAt, &recurringTransaction.UpdatedAt); err != nil {
			return nil, err
		}
		recurringTransactions = append(recurringTransactions, recurringTransaction)
	}
	return recurringTransactions, nil
}

func GetRecurringTransactions(db interfaces.SqlExecutor) ([]models.RecurringTransaction, error) {
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurringTransactions []models.RecurringTransaction
	for rows.Next() {
		var recurringTransaction models.RecurringTransaction
		if err := rows.Scan(&recurringTransaction.ID, &recurringTransaction.UserID, &recurringTransaction.AccountID, &recurringTransaction.CategoryID, &recurringTransaction.BudgetID, &recurringTransaction.Description, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.Note, &recurringTransaction.RecurringFrequency, &recurringTransaction.RecurringDate, &recurringTransaction.CreatedAt, &recurringTransaction.UpdatedAt); err != nil {
			return nil, err
		}
		recurringTransactions = append(recurringTransactions, recurringTransaction)
	}
	return recurringTransactions, nil
}

func GetRecurringTransactionByID(id uuid.UUID, db interfaces.SqlExecutor) (*models.RecurringTransaction, error) {
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions WHERE id = $1"
	row := db.QueryRow(query, id)

	var recurringTransaction models.RecurringTransaction
	if err := row.Scan(&recurringTransaction.ID, &recurringTransaction.UserID, &recurringTransaction.AccountID, &recurringTransaction.CategoryID, &recurringTransaction.BudgetID, &recurringTransaction.Description, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.Note, &recurringTransaction.RecurringFrequency, &recurringTransaction.RecurringDate, &recurringTransaction.CreatedAt, &recurringTransaction.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &recurringTransaction, nil
}

func UpdateRecurringTransaction(recurringTransaction *models.RecurringTransaction, db interfaces.SqlExecutor) error {
	query := "UPDATE recurring_transactions SET account_id = $1, category_id = $2, budget_id = $3, description = $4, amount = $5, type = $6, note = $7, recurring_frequency = $8, recurring_date = $9, updated_at = $10 WHERE id = $11"
	_, err := db.Exec(query, recurringTransaction.AccountID, recurringTransaction.CategoryID, recurringTransaction.BudgetID, recurringTransaction.Description, recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.Note, recurringTransaction.RecurringFrequency, recurringTransaction.RecurringDate, recurringTransaction.UpdatedAt, recurringTransaction.ID)
	return err
}

func DeleteRecurringTransaction(id uuid.UUID, db interfaces.SqlExecutor) error {
	query := "DELETE FROM recurring_transactions WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
