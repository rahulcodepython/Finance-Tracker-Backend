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
	Frequency   string          `json:"frequency"`
	StartDate   time.Time       `json:"startDate"`
	EndDate     time.Time       `json:"endDate"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

var RecurringTransactionColumns = "id, user_id, account_id, category_id, description, amount, type, frequency, start_date, end_date, is_active, created_at, updated_at"
