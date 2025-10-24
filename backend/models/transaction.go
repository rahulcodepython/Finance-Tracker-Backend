package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TransactionType defines the set of possible transaction types.
type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"userId"`
	AccountID       uuid.UUID       `json:"accountId"`
	CategoryID      uuid.UUID       `json:"categoryId"`
	BudgetID        uuid.NullUUID   `json:"budgetId,omitempty"`
	Description     string          `json:"description"`
	Amount          float64         `json:"amount"` // See note on NUMERIC type above
	Type            TransactionType `json:"type"`
	TransactionDate time.Time       `json:"transactionDate"`
	Note            sql.NullString  `json:"note,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

var TransactionColumns = "id, user_id, account_id, category_id, budget_id, description, amount, type, transaction_date, note, created_at, updated_at"
