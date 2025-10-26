// Package models defines the data structures used in the application.
package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// RecurringFrequency defines the set of possible frequencies for recurring transactions.
type RecurringFrequency string

// Constants for the different recurring frequencies.
const (
	Monthly RecurringFrequency = "monthly"
	Yearly  RecurringFrequency = "yearly"
)

// RecurringTransaction represents a recurring transaction in the `recurring_transactions` table.
type RecurringTransaction struct {
	ID                 uuid.UUID          `json:"id"`                 // Unique identifier for the recurring transaction
	UserID             uuid.UUID          `json:"userId"`             // ID of the user who owns the transaction
	AccountID          uuid.UUID          `json:"accountId"`          // ID of the account associated with the transaction
	CategoryID         uuid.UUID          `json:"categoryId"`         // ID of the category associated with the transaction
	BudgetID           uuid.NullUUID      `json:"budgetId,omitempty"`   // ID of the budget associated with the transaction (optional)
	Description        string             `json:"description"`        // Description of the transaction
	Amount             float64            `json:"amount"`             // Amount of the transaction
	Type               TransactionType    `json:"type"`               // Type of the transaction (e.g., income or expense)
	Note               sql.NullString     `json:"note,omitempty"`      // Additional notes for the transaction (optional)
	RecurringFrequency RecurringFrequency `json:"recurringFrequency"` // Frequency of the recurring transaction
	RecurringDate      int                `json:"recurringDate"`      // Day of the month or year for the recurring transaction
	CreatedAt          time.Time          `json:"createdAt"`          // Timestamp of when the transaction was created
	UpdatedAt          time.Time          `json:"updatedAt"`          // Timestamp of when the transaction was last updated
}

// RecurringTransactionColumns is a comma-separated string of all recurring transaction columns, useful for SQL queries.
var RecurringTransactionColumns = "id, user_id, account_id, category_id, budget_id, description, amount, type, note, recurring_frequency, recurring_date, created_at, updated_at"