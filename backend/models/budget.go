// Package models defines the data structures used in the application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Budget represents a user-defined budget in the `budgets` table.
type Budget struct {
	ID        uuid.UUID `json:"id"`        // Unique identifier for the budget
	UserID    uuid.UUID `json:"userId"`    // ID of the user who owns the budget
	Name      string    `json:"name"`      // Name of the budget
	Amount    float64   `json:"amount"`    // The allocated amount for the budget
	CreatedAt time.Time `json:"createdAt"` // Timestamp of when the budget was created
	UpdatedAt time.Time `json:"updatedAt"` // Timestamp of when the budget was last updated
}

// BudgetColumns is a comma-separated string of all budget columns, useful for SQL queries.
var BudgetColumns = "id, user_id, name, amount, created_at, updated_at"