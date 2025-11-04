package serializers

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateRecurringTransactionInput defines the input structure for creating a new recurring transaction.
type RecurringTransactionInput struct {
	AccountID          string                    `json:"accountId"`
	CategoryID         string                    `json:"categoryId"`
	BudgetID           string                    `json:"budgetId"`
	Description        string                    `json:"description"`
	Amount             float64                   `json:"amount"`
	Note               string                    `json:"note"`
	RecurringFrequency models.RecurringFrequency `json:"recurringFrequency"`
	RecurringDate      int                       `json:"recurringDate"`
}

type RecurringTransactionResponse struct {
	ID                 uuid.UUID                 `json:"id"`
	UserID             uuid.UUID                 `json:"userId"`
	Description        string                    `json:"description"`
	Amount             float64                   `json:"amount"`
	Type               models.TransactionType    `json:"type"`
	Note               sql.NullString            `json:"note,omitempty"`
	RecurringFrequency models.RecurringFrequency `json:"recurringFrequency"`
	RecurringDate      int                       `json:"recurringDate"`
	CreatedAt          time.Time                 `json:"createdAt"`
	UpdatedAt          time.Time                 `json:"updatedAt"`

	Account  AccountRef  `json:"account"`
	Category CategoryRef `json:"category"`
	Budget   BudgetRef   `json:"budget"`
}
