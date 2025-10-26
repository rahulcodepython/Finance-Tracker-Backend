// Package models defines the data structures used in the application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// AccountType defines the set of possible account types.
type AccountType string

// Constants for the different account types.
const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeCash       AccountType = "cash"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeLoan       AccountType = "loan"
	AccountTypeUPI        AccountType = "upi"
)

// Account represents a user's financial account in the `accounts` table.
type Account struct {
	ID        uuid.UUID   `json:"id"`        // Unique identifier for the account
	UserID    uuid.UUID   `json:"userId"`    // ID of the user who owns the account
	Name      string      `json:"name"`      // Name of the account
	Type      AccountType `json:"type"`      // Type of the account
	Balance   float64     `json:"balance"`   // Current balance of the account
	IsActive  bool        `json:"isActive"`  // Whether the account is active
	CreatedAt time.Time   `json:"createdAt"` // Timestamp of when the account was created
	UpdatedAt time.Time   `json:"updatedAt"` // Timestamp of when the account was last updated
}

// AccountColumns is a comma-separated string of all account columns, useful for SQL queries.
var AccountColumns = "id, user_id, name, type, balance, is_active, created_at, updated_at"