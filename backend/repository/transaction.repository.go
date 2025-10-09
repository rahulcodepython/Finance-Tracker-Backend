package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateTransaction(transaction *models.Transaction, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO transactions (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", models.TransactionColumns)
	_, err := db.Exec(query, transaction.ID, transaction.UserID, transaction.AccountID, transaction.CategoryID, transaction.Description, transaction.Amount, transaction.Type, transaction.TransactionDate, transaction.Note, transaction.CreatedAt, transaction.UpdatedAt)
	return err
}

func GetTransactionsByUserID(userID uuid.UUID, db *sql.DB) ([]models.Transaction, error) {
	query := "SELECT * FROM transactions WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID, &transaction.Description, &transaction.Amount, &transaction.Type, &transaction.TransactionDate, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionByID(id uuid.UUID, db *sql.DB) (*models.Transaction, error) {
	query := "SELECT * FROM transactions WHERE id = $1"
	row := db.QueryRow(query, id)

	var transaction models.Transaction
	if err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID, &transaction.Description, &transaction.Amount, &transaction.Type, &transaction.TransactionDate, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &transaction, nil
}

func UpdateTransaction(transaction *models.Transaction, db *sql.DB) error {
	query := "UPDATE transactions SET account_id = $1, category_id = $2, description = $3, amount = $4, type = $5, transaction_date = $6, note = $7, updated_at = $8 WHERE id = $9"
	_, err := db.Exec(query, transaction.AccountID, transaction.CategoryID, transaction.Description, transaction.Amount, transaction.Type, transaction.TransactionDate, transaction.Note, transaction.UpdatedAt, transaction.ID)
	return err
}

func DeleteTransaction(id uuid.UUID, db *sql.DB) error {
	query := "DELETE FROM transactions WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}

func GetTransactionsByUserIDWithFilters(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, startDate string, endDate string, db *sql.DB) ([]models.Transaction, error) {
	var query strings.Builder
	query.WriteString("SELECT * FROM transactions WHERE user_id = $1")

	args := []interface{}{userID}
	argCount := 2

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

	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit))

	rows, err := db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID, &transaction.Description, &transaction.Amount, &transaction.Type, &transaction.TransactionDate, &transaction.Note, &transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetAggregateDataByUserID(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	var totalIncome float64
	var totalExpenses float64

	var query strings.Builder
	query.WriteString("SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income, COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expenses FROM transactions WHERE user_id = $1")

	args := []interface{}{userID}
	argCount := 2

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

	row := db.QueryRow(query.String(), args...)
	if err := row.Scan(&totalIncome, &totalExpenses); err != nil {
		return nil, err
	}

	netIncome := totalIncome - totalExpenses

	return map[string]interface{}{
		"totalIncome":   totalIncome,
		"totalExpenses": totalExpenses,
		"netIncome":     netIncome,
	}, nil
}

func GetSpendingByCategory(userID uuid.UUID, db *sql.DB) ([]map[string]interface{}, error) {
	query := "SELECT c.name as category, sum(t.amount) as amount FROM transactions t JOIN categories c ON c.id = t.category_id WHERE t.user_id = $1 AND t.type = 'expense' GROUP BY c.name"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
