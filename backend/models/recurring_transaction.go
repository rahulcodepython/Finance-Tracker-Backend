package models

import (
	"time"

	"github.com/google/uuid"
)

// RecurringTransaction corresponds to the `recurring_transactions` table.

type RecurringTransaction struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"userId"`
	AccountID   uuid.UUID       `json:"accountId"`
	CategoryID  uuid.NullUUID   `json:"categoryId"`
	Description string          `json:"description"`
	Amount      float64         `json:"amount"`
	Type        TransactionType `json:"type"`
	RecurringFrequency   RecurringFrequency          `json:"recurringFrequency"`
	RecurringDate      int       `json:"recurringDate"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

var RecurringTransactionColumns = "id, user_id, account_id, category_id, description, amount, type, recurring_frequency, recurring_date, created_at, updated_at"
