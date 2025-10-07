package models

import (
	"time"

	"github.com/google/uuid"
)

// Budget corresponds to the `budgets` table.
type Budget struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	CategoryID uuid.UUID `json:"categoryId"`
	Amount     float64   `json:"amount"`
	Month      time.Time `json:"month"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

var BudgetColumns = "id, user_id, category_id, amount, month, created_at, updated_at"
