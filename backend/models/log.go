// Package models defines the data structures used in the application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Log represents a user activity log in the `logs` table.
type Log struct {
	ID        uuid.UUID `json:"id"`        // Unique identifier for the log entry
	UserID    uuid.UUID `json:"userId"`    // ID of the user who performed the action
	Message   string    `json:"message"`   // The log message
	CreatedAt time.Time `json:"createdAt"` // Timestamp of when the log was created
}

// LogColumns is a comma-separated string of all log columns, useful for SQL queries.
const LogColumns = "id, user_id, message, created_at"
