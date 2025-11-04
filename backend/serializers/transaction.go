// Package serializers defines the data structures used in the application.
package serializers

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateTransactionInput defines the input structure for creating a new transaction.
type TransactionInput struct {
	AccountID   string  `json:"accountId"`
	CategoryID  string  `json:"categoryId"`
	BudgetID    string  `json:"budgetId"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Note        string  `json:"note"`
}

type TransactionResponse struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"userId"`
	Description     string                 `json:"description"`
	Amount          float64                `json:"amount"`
	Type            models.TransactionType `json:"type"`
	TransactionDate time.Time              `json:"transactionDate"`
	Note            sql.NullString         `json:"note,omitempty"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`

	Account  AccountRef  `json:"account"`
	Category CategoryRef `json:"category"`
	Budget   BudgetRef   `json:"budget"`
}

type AccountRef struct {
	UUID uuid.UUID `json:"UUID"`
	Name string    `json:"name"`
}

type CategoryRef struct {
	UUID uuid.UUID `json:"UUID"`
	Name string    `json:"name"`
}

type BudgetRef struct {
	UUID  uuid.UUID      `json:"UUID"`
	Name  sql.NullString `json:"name,omitempty"`
	Valid bool           `json:"Valid"`
}
