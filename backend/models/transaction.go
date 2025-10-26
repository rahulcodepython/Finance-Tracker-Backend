// Package models defines the data structures used in the application.
package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TransactionType defines the set of possible transaction types.
type TransactionType string

// Constants for the different transaction types.
const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

// Transaction represents a financial transaction in the `transactions` table.
type Transaction struct {
	ID              uuid.UUID       `json:"id"`                 // Unique identifier for the transaction
	UserID          uuid.UUID       `json:"userId"`             // ID of the user who owns the transaction
	AccountID       uuid.UUID       `json:"accountId"`          // ID of the account associated with the transaction
	CategoryID      uuid.UUID       `json:"categoryId"`         // ID of the category associated with the transaction
	BudgetID        uuid.NullUUID   `json:"budgetId,omitempty"`   // ID of the budget associated with the transaction (optional)
	Description     string          `json:"description"`        // Description of the transaction
	Amount          float64         `json:"amount"`             // Amount of the transaction
	Type            TransactionType `json:"type"`               // Type of the transaction (e.g., income or expense)
	TransactionDate time.Time       `json:"transactionDate"`    // Date of the transaction
	Note            sql.NullString  `json:"note,omitempty"`      // Additional notes for the transaction (optional)
	CreatedAt       time.Time       `json:"createdAt"`          // Timestamp of when the transaction was created
	UpdatedAt       time.Time       `json:"updatedAt"`          // Timestamp of when the transaction was last updated
}

// TransactionColumns is a comma-separated string of all transaction columns, useful for SQL queries.
var TransactionColumns = "id, user_id, account_id, category_id, budget_id, description, amount, type, transaction_date, note, created_at, updated_at"