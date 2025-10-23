package models

import (
	"time"

	"github.com/google/uuid"
)

// Budget corresponds to the `budgets` table.
type Budget struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var BudgetColumns = "id, user_id, name, amount, created_at, updated_at"
