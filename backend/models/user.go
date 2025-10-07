package models

import (
	"time"

	"github.com/google/uuid"
)

// AuthProvider defines the set of possible authentication providers.
type AuthProvider string

const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
)

// User corresponds to the `users` table.
type User struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Password  string       `json:"-"` // Omitted from JSON responses for security
	Provider  AuthProvider `json:"provider"`
	CreatedAt time.Time    `json:"createdAt"`
}

var UserColumns = "id, name, email, password, provider, created_at"
