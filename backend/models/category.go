// Package models defines the data structures used in the application.
package models

import "github.com/google/uuid"

// Category represents a transaction category in the `categories` table.
type Category struct {
	ID   uuid.UUID       `json:"id"`   // Unique identifier for the category
	Name string          `json:"name"` // Name of the category
	Type TransactionType `json:"type"` // Type of the category (e.g., income or expense)
}

// CategoryColumns is a comma-separated string of all category columns, useful for SQL queries.
var CategoryColumns = "id, name, type"