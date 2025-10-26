// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateTransaction inserts a new transaction record into the database.
func CreateTransaction(transaction *models.Transaction, db interfaces.SqlExecutor) error {
	// Construct the SQL query for insertion.
	query := fmt.Sprintf("INSERT INTO transactions (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", models.TransactionColumns)
	// Execute the query with the transaction data.
	_, err := db.Exec(query, transaction.ID, transaction.UserID, transaction.AccountID, transaction.CategoryID, transaction.BudgetID, transaction.Description, transaction.Amount, transaction.Type, transaction.TransactionDate, transaction.Note, transaction.CreatedAt, transaction.UpdatedAt)
	return err
}

// GetTransactionsByUserID retrieves all transactions for a given user ID from the database.
func GetTransactionsByUserID(userID uuid.UUID, db interfaces.SqlExecutor) ([]models.Transaction, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM transactions WHERE user_id = $1"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into Transaction structs.
	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID, &transaction.BudgetID, &transaction.Description, &transaction.Amount, &transaction.Type, &transaction.TransactionDate, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// GetTransactionByID retrieves a single transaction by its ID from the database.
func GetTransactionByID(id uuid.UUID, db interfaces.SqlExecutor) (*models.Transaction, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM transactions WHERE id = $1"
	// Execute the query.
	row := db.QueryRow(query, id)

	// Scan the row into a Transaction struct.
	var transaction models.Transaction
	if err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID, &transaction.BudgetID, &transaction.Description, &transaction.Amount, &transaction.Type, &transaction.TransactionDate, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, nil
}

// UpdateTransaction updates an existing transaction record in the database.
func UpdateTransaction(transaction *models.Transaction, db interfaces.SqlExecutor) error {
	// Construct the SQL query for update.
	query := "UPDATE transactions SET account_id = $1, category_id = $2, budget_id = $3, description = $4, amount = $5, type = $6, transaction_date = $7, note = $8, updated_at = $9 WHERE id = $10"
	// Execute the query with the updated transaction data.
	_, err := db.Exec(query, transaction.AccountID, transaction.CategoryID, transaction.BudgetID, transaction.Description, transaction.Amount, transaction.Type, transaction.TransactionDate, transaction.Note, transaction.UpdatedAt, transaction.ID)
	return err
}

// DeleteTransaction deletes a transaction record from the database by its ID.
func DeleteTransaction(id uuid.UUID, db interfaces.SqlExecutor) error {
	// Construct the SQL query for deletion.
	query := "DELETE FROM transactions WHERE id = $1"
	// Execute the query.
	_, err := db.Exec(query, id)
	return err
}

// GetTransactionsByUserIDWithFilters retrieves transactions for a given user ID with optional filters and pagination.
func GetTransactionsByUserIDWithFilters(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, budgetID string, startDate string, endDate string, db interfaces.SqlExecutor) ([]models.Transaction, error) {
	// Use a strings.Builder to efficiently construct the SQL query.
	var query strings.Builder
	query.WriteString("SELECT id, user_id, account_id, category_id, budget_id, description, amount, type, transaction_date, note, created_at, updated_at FROM transactions WHERE user_id = $1")

	// Create a slice to hold the query arguments.
	args := []interface{}{userID}
	argCount := 2

	// Add filters to the query based on the provided parameters.
	if description != "" {
		query.WriteString(fmt.Sprintf(" AND description LIKE $%d", argCount))
		args = append(args, "%"+description+"%")
		argCount++
	}

	if categoryID != "" {
		query.WriteString(fmt.Sprintf(" AND category_id = $%d", argCount))
		args = append(args, categoryID)
		argCount++
	}

	if accountID != "" {
		query.WriteString(fmt.Sprintf(" AND account_id = $%d", argCount))
		args = append(args, accountID)
		argCount++
	}

	if budgetID != "" {
		query.WriteString(fmt.Sprintf(" AND budget_id = $%d", argCount))
		args = append(args, budgetID)
		argCount++
	}

	if startDate != "" {
		query.WriteString(fmt.Sprintf(" AND transaction_date >= $%d", argCount))
		args = append(args, startDate)
		argCount++
	}

	if endDate != "" {
		query.WriteString(fmt.Sprintf(" AND transaction_date <= $%d", argCount))
		args = append(args, endDate)
		argCount++
	}

	// Add pagination to the query.
	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit))

	// Execute the query.
	rows, err := db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into Transaction structs.
	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID, &transaction.BudgetID, &transaction.Description, &transaction.Amount, &transaction.Type, &transaction.TransactionDate, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// GetAggregateDataByUserID retrieves aggregate transaction data for a given user ID and date range.
func GetAggregateDataByUserID(userID uuid.UUID, startDate string, endDate string, db interfaces.SqlExecutor) (map[string]interface{}, error) {
	var totalIncome float64
	var totalExpenses float64

	// Use a strings.Builder to efficiently construct the SQL query.
	var query strings.Builder
	query.WriteString("SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income, COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expenses FROM transactions WHERE user_id = $1")

	// Create a slice to hold the query arguments.
	args := []interface{}{userID}
	argCount := 2

	// Add date range filtering to the query if start and end dates are provided.
	if startDate != "" {
		query.WriteString(fmt.Sprintf(" AND transaction_date >= $%d", argCount))
		args = append(args, startDate)
		argCount++
	}

	if endDate != "" {
		query.WriteString(fmt.Sprintf(" AND transaction_date <= $%d", argCount))
		args = append(args, endDate)
		argCount++
	}

	// Execute the query and scan the results.
	row := db.QueryRow(query.String(), args...)
	if err := row.Scan(&totalIncome, &totalExpenses); err != nil {
		return nil, err
	}

	// Calculate the net income.
	netIncome := totalIncome - totalExpenses

	// Return the aggregate data in a map.
	return map[string]interface{}{
		"totalIncome":   totalIncome,
		"totalExpenses": totalExpenses,
		"netIncome":     netIncome,
	}, nil
}

// GetSpendingByCategory retrieves the total spending by category for a given user ID.
func GetSpendingByCategory(userID uuid.UUID, db interfaces.SqlExecutor) ([]map[string]interface{}, error) {
	// Construct the SQL query for selection.
	query := "SELECT c.name as category, sum(t.amount) as amount FROM transactions t JOIN categories c ON c.id = t.category_id WHERE t.user_id = $1 AND t.type = 'expense' GROUP BY c.name"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into a slice of maps.
	var result []map[string]interface{}
	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{"category": category, "amount": amount})
	}

	return result, nil
}

// GetEarningByCategory retrieves the total earnings by category for a given user ID.
func GetEarningByCategory(userID uuid.UUID, db interfaces.SqlExecutor) ([]map[string]interface{}, error) {
	// Construct the SQL query for selection.
	query := "SELECT c.name as category, sum(t.amount) as amount FROM transactions t JOIN categories c ON c.id = t.category_id WHERE t.user_id = $1 AND t.type = 'income' GROUP BY c.name"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into a slice of maps.
	var result []map[string]interface{}
	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{"category": category, "amount": amount})
	}

	return result, nil
}