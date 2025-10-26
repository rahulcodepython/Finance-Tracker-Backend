// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateRecurringTransaction inserts a new recurring transaction record into the database.
func CreateRecurringTransaction(recurringTransaction *models.RecurringTransaction, db interfaces.SqlExecutor) error {
	// Construct the SQL query for insertion.
	query := fmt.Sprintf("INSERT INTO recurring_transactions (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", models.RecurringTransactionColumns)
	// Execute the query with the recurring transaction data.
	_, err := db.Exec(query, recurringTransaction.ID, recurringTransaction.UserID, recurringTransaction.AccountID, recurringTransaction.CategoryID, recurringTransaction.BudgetID, recurringTransaction.Description, recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.Note, recurringTransaction.RecurringFrequency, recurringTransaction.RecurringDate, recurringTransaction.CreatedAt, recurringTransaction.UpdatedAt)
	return err
}

// GetRecurringTransactionsByUserID retrieves all recurring transactions for a given user ID from the database.
func GetRecurringTransactionsByUserID(userID uuid.UUID, db interfaces.SqlExecutor) ([]models.RecurringTransaction, error) {
	// Construct the SQL query for selection.
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions WHERE user_id = $1"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into RecurringTransaction structs.
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

// GetRecurringTransactions retrieves all recurring transactions from the database.
func GetRecurringTransactions(db interfaces.SqlExecutor) ([]models.RecurringTransaction, error) {
	// Construct the SQL query for selection.
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions"
	// Execute the query.
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into RecurringTransaction structs.
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

// GetRecurringTransactionByID retrieves a single recurring transaction by its ID from the database.
func GetRecurringTransactionByID(id uuid.UUID, db interfaces.SqlExecutor) (*models.RecurringTransaction, error) {
	// Construct the SQL query for selection.
	query := "SELECT " + models.RecurringTransactionColumns + " FROM recurring_transactions WHERE id = $1"
	// Execute the query.
	row := db.QueryRow(query, id)

	// Scan the row into a RecurringTransaction struct.
	var recurringTransaction models.RecurringTransaction
	if err := row.Scan(&recurringTransaction.ID, &recurringTransaction.UserID, &recurringTransaction.AccountID, &recurringTransaction.CategoryID, &recurringTransaction.BudgetID, &recurringTransaction.Description, &recurringTransaction.Amount, &recurringTransaction.Type, &recurringTransaction.Note, &recurringTransaction.RecurringFrequency, &recurringTransaction.RecurringDate, &recurringTransaction.CreatedAt, &recurringTransaction.UpdatedAt); err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &recurringTransaction, nil
}

// UpdateRecurringTransaction updates an existing recurring transaction record in the database.
func UpdateRecurringTransaction(recurringTransaction *models.RecurringTransaction, db interfaces.SqlExecutor) error {
	// Construct the SQL query for update.
	query := "UPDATE recurring_transactions SET account_id = $1, category_id = $2, budget_id = $3, description = $4, amount = $5, type = $6, note = $7, recurring_frequency = $8, recurring_date = $9, updated_at = $10 WHERE id = $11"
	// Execute the query with the updated recurring transaction data.
	_, err := db.Exec(query, recurringTransaction.AccountID, recurringTransaction.CategoryID, recurringTransaction.BudgetID, recurringTransaction.Description, recurringTransaction.Amount, recurringTransaction.Type, recurringTransaction.Note, recurringTransaction.RecurringFrequency, recurringTransaction.RecurringDate, recurringTransaction.UpdatedAt, recurringTransaction.ID)
	return err
}

// DeleteRecurringTransaction deletes a recurring transaction record from the database by its ID.
func DeleteRecurringTransaction(id uuid.UUID, db interfaces.SqlExecutor) error {
	// Construct the SQL query for deletion.
	query := "DELETE FROM recurring_transactions WHERE id = $1"
	// Execute the query.
	_, err := db.Exec(query, id)
	return err
}