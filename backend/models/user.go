// Package models defines the data structures used in the application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider defines the set of possible authentication providers.
type AuthProvider string

// Constants for the different authentication providers.
const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
)

// User represents a user of the application in the `users` table.
type User struct {
	ID        uuid.UUID    `json:"id"`        // Unique identifier for the user
	Name      string       `json:"name"`      // Name of the user
	Email     string       `json:"email"`      // Email address of the user (must be unique)
	Password  string       `json:"-"`         // User's password (omitted from JSON responses for security)
	Provider  AuthProvider `json:"provider"`  // Authentication provider used by the user
	CreatedAt time.Time    `json:"createdAt"` // Timestamp of when the user was created
}

// UserColumns is a comma-separated string of all user columns, useful for SQL queries.
var UserColumns = "id, name, email, password, provider, created_at"