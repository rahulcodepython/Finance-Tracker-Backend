package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// RecurringTransaction corresponds to the `recurring_transactions` table.
type RecurringFrequency string

const (
	Monthly RecurringFrequency = "monthly"
	Yearly  RecurringFrequency = "yearly"
)

type RecurringTransaction struct {
	ID                 uuid.UUID          `json:"id"`
	UserID             uuid.UUID          `json:"userId"`
	AccountID          uuid.UUID          `json:"accountId"`
	CategoryID         uuid.UUID          `json:"categoryId"`
	BudgetID           uuid.NullUUID      `json:"budgetId,omitempty"`
	Description        string             `json:"description"`
	Amount             float64            `json:"amount"`
	Type               TransactionType    `json:"type"`
	Note               sql.NullString     `json:"note,omitempty"`
	RecurringFrequency RecurringFrequency `json:"recurringFrequency"`
	RecurringDate      int                `json:"recurringDate"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

var RecurringTransactionColumns = "id, user_id, account_id, category_id, budget_id, description, amount, type, note, recurring_frequency, recurring_date, created_at, updated_at"
