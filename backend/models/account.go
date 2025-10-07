package models

import (
	"time"

	"github.com/google/uuid"
)

// AccountType defines the set of possible account types.
type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeCash       AccountType = "cash"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeLoan       AccountType = "loan"
	AccountTypeUPI        AccountType = "upi"
)

// Account corresponds to the `accounts` table.
type Account struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"userId"`
	Name      string      `json:"name"`
	Type      AccountType `json:"type"`
	Balance   float64     `json:"balance"` // See note on NUMERIC type above
	IsActive  bool        `json:"isActive"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

var AccountColumns = "id, user_id, name, type, balance, is_active, created_at, updated_at"
