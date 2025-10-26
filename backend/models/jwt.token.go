// Package models defines the data structures used in the application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// JwtToken represents a JWT token stored in the `jwt_tokens` table.
// This is used to keep track of active tokens and facilitate token revocation.
type JwtToken struct {
	ID        uuid.UUID `json:"id"`        // Unique identifier for the token record
	UserID    uuid.UUID `json:"userId"`    // ID of the user to whom the token was issued
	Token     string    `json:"token"`     // The JWT token string
	ExpiresAt time.Time `json:"expiresAt"` // Timestamp of when the token expires
	CreatedAt time.Time `json:"createdAt"` // Timestamp of when the token was created
}

// JwtTokenColumns is a comma-separated string of all JWT token columns, useful for SQL queries.
var JwtTokenColumns = "id, user_id, token, expires_at, created_at"