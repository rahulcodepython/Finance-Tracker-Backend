// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/serializers"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
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
func GetTransactionsByUserIDWithFilters(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, budgetID string, startDate string, endDate string, db interfaces.SqlExecutor) ([]serializers.TransactionResponse, error) {
	// Use a strings.Builder to efficiently construct the SQL query.
	var query strings.Builder
	query.WriteString(`
	SELECT 
        t.id, t.user_id, t.description, t.amount, t.type, t.transaction_date,
        t.note, t.created_at, t.updated_at,
        a.id, a.name,
        c.id, c.name,
        b.id, b.name
    FROM transactions t
    INNER JOIN accounts a ON t.account_id = a.id
    INNER JOIN categories c ON t.category_id = c.id
    LEFT JOIN budgets b ON t.budget_id = b.id
    WHERE t.user_id = $1

	`)

	// Create a slice to hold the query arguments.
	args := []interface{}{userID}
	argCount := 2

	// Add filters to the query based on the provided parameters.
	if description != "" {
		query.WriteString(fmt.Sprintf(" AND t.description LIKE $%d", argCount))
		args = append(args, "%"+description+"%")
		argCount++
	}

	if categoryID != "" {
		query.WriteString(fmt.Sprintf(" AND t.category_id = $%d", argCount))
		args = append(args, categoryID)
		argCount++
	}

	if accountID != "" {
		query.WriteString(fmt.Sprintf(" AND t.account_id = $%d", argCount))
		args = append(args, accountID)
		argCount++
	}

	if budgetID != "" {
		query.WriteString(fmt.Sprintf(" AND t.budget_id = $%d", argCount))
		args = append(args, budgetID)
		argCount++
	}

	if startDate != "" {
		query.WriteString(fmt.Sprintf(" AND t.transaction_date >= $%d", argCount))
		args = append(args, startDate)
		argCount++
	}

	if endDate != "" {
		query.WriteString(fmt.Sprintf(" AND t.transaction_date <= $%d", argCount))
		args = append(args, endDate)
		argCount++
	}

	// Add sorting to the query
	query.WriteString(" ORDER BY t.transaction_date DESC")

	// Add pagination to the query.
	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit))

	// Execute the query.
	rows, err := db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []serializers.TransactionResponse

	for rows.Next() {
		var tr serializers.TransactionResponse
		err := rows.Scan(
			&tr.ID, &tr.UserID, &tr.Description, &tr.Amount, &tr.Type, &tr.TransactionDate,
			&tr.Note, &tr.CreatedAt, &tr.UpdatedAt,
			&tr.Account.UUID, &tr.Account.Name,
			&tr.Category.UUID, &tr.Category.Name,
			&tr.Budget.UUID, &tr.Budget.Name,
		)
		if err != nil {
			return nil, err
		}

		// budget validity check
		tr.Budget.Valid = tr.Budget.UUID != uuid.Nil

		transactions = append(transactions, tr)
	}

	return transactions, nil
}

// GetAggregateDataByUserID retrieves aggregate transaction data for a given user ID and date range.
func GetAggregateDataByUserID(userID uuid.UUID, startDate string, endDate string, db interfaces.SqlExecutor) (*serializers.DashboardSummary, error) {
	var dashboardSummer serializers.DashboardSummary

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
	} else {
		query.WriteString(fmt.Sprintf(" AND transaction_date >= $%d", argCount))
		args = append(args, time.Now().In(utils.LOC).AddDate(0, 0, -30))
		argCount++
	}

	if endDate != "" {
		query.WriteString(fmt.Sprintf(" AND transaction_date <= $%d", argCount))
		args = append(args, endDate)
		argCount++
	} else {
		query.WriteString(fmt.Sprintf(" AND transaction_date >= $%d", argCount))
		args = append(args, time.Now().In(utils.LOC))
		argCount++
	}

	// Execute the query and scan the results.
	row := db.QueryRow(query.String(), args...)
	if err := row.Scan(&dashboardSummer.MonthlyIncome, &dashboardSummer.MonthlyExpenses); err != nil {
		return nil, err
	}

	// Calculate the net income.
	dashboardSummer.MonthlySavings = dashboardSummer.MonthlyIncome - dashboardSummer.MonthlyExpenses

	return &dashboardSummer, nil
}

// GetSpendingByCategory retrieves the total spending by category for a given user ID.
func GetSpendingByCategory(userID uuid.UUID, db interfaces.SqlExecutor) (*[]serializers.CategoryAggregate, error) {
	// Construct the SQL query for selection.
	query := "SELECT c.name as category, sum(t.amount) as amount FROM transactions t JOIN categories c ON c.id = t.category_id WHERE t.user_id = $1 AND t.type = 'expense' GROUP BY c.name"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into a slice of maps.
	var result []serializers.CategoryAggregate
	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			return nil, err
		}

		var value serializers.CategoryAggregate
		value.Category = category
		value.Amount = amount

		result = append(result, value)
	}

	return &result, nil
}

// GetEarningByCategory retrieves the total earnings by category for a given user ID.
func GetEarningByCategory(userID uuid.UUID, db interfaces.SqlExecutor) (*[]serializers.CategoryAggregate, error) {
	// Construct the SQL query for selection.
	query := "SELECT c.name as category, sum(t.amount) as amount FROM transactions t JOIN categories c ON c.id = t.category_id WHERE t.user_id = $1 AND t.type = 'income' GROUP BY c.name"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into a slice of maps.
	var result []serializers.CategoryAggregate
	for rows.Next() {
		var category string
		var amount float64
		if err := rows.Scan(&category, &amount); err != nil {
			return nil, err
		}

		var value serializers.CategoryAggregate
		value.Category = category
		value.Amount = amount

		result = append(result, value)
	}

	return &result, nil
}

func GetIncomeExpense(userID uuid.UUID, db interfaces.SqlExecutor) (*serializers.IncomeExpenseAggregate, error) {
	var incomeExpenseAggregate serializers.IncomeExpenseAggregate

	query := ` SELECT SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) AS total_income, SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) AS total_expense FROM transactions WHERE user_id = $1;`

	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&incomeExpenseAggregate.Income, &incomeExpenseAggregate.Expense); err != nil {
			return nil, err
		}
	}

	return &incomeExpenseAggregate, nil
}
