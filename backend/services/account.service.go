// Package services provides business logic for the application.
package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateAccount creates a new financial account for a user.
func CreateAccount(userID uuid.UUID, name string, accountType models.AccountType, balance float64, db *sql.DB) (*models.Account, error) {
	// Create a new Account model.
	account := &models.Account{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Type:      accountType,
		Balance:   balance,
		IsActive:  true,
		CreatedAt: time.Now().In(utils.LOC),
		UpdatedAt: time.Now().In(utils.LOC),
	}

	// Create the account in the database.
	err := repository.CreateAccount(account, db)
	if err != nil {
		return nil, err
	}

	// Create a log entry for the account creation.
	go CreateLog(userID, fmt.Sprintf("New account '%s' created", account.Name), db)

	return account, nil
}

// GetAccounts retrieves all financial accounts for a user.
func GetAccounts(userID uuid.UUID, db *sql.DB) ([]models.Account, error) {
	return repository.GetAccountsByUserID(userID, db)
}

// CheckAccountExistsById checks if an account exists by its ID.
func CheckAccountExistsById(id uuid.UUID, db *sql.DB) (bool, error) {
	// Get the account from the database.
	account, err := repository.GetAccountByID(id, db)
	if err != nil {
		// If no rows were found, the account does not exist.
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	// If the account is not nil, it exists.
	if account != nil {
		return true, nil
	}

	return false, nil
}

// UpdateAccount updates an existing financial account.
func UpdateAccount(id uuid.UUID, name string, accountType models.AccountType, isActive bool, db *sql.DB) (*models.Account, error) {
	// Get the account from the database.
	account, err := repository.GetAccountByID(id, db)
	if err != nil {
		return nil, err
	}

	// If the account does not exist, return an error.
	if account == nil {
		return nil, sql.ErrNoRows
	}

	// Update the account fields.
	account.Name = name
	account.Type = accountType
	account.IsActive = isActive
	account.UpdatedAt = time.Now().In(utils.LOC)

	// Update the account in the database.
	err = repository.UpdateAccount(account, db)
	if err != nil {
		return nil, err
	}

	// Create a log entry for the account update.
	go CreateLog(account.UserID, fmt.Sprintf("Account '%s' updated", account.Name), db)

	return account, nil
}

// DeleteAccount deletes a financial account.
func DeleteAccount(id uuid.UUID, db *sql.DB) error {
	// Get the account from the database.
	account, err := repository.GetAccountByID(id, db)
	if err != nil {
		return err
	}

	// If the account does not exist, return an error.
	if account == nil {
		return sql.ErrNoRows
	}

	// Delete the account from the database.
	err = repository.DeleteAccount(id, db)
	if err != nil {
		return err
	}

	// Create a log entry for the account deletion.
	go CreateLog(account.UserID, fmt.Sprintf("Account '%s' removed", account.Name), db)

	return nil
}

// GetTotalBalance calculates the total balance of all active accounts for a user.
func GetTotalBalance(userID uuid.UUID, db *sql.DB) (float64, error) {
	// Get all accounts for the user.
	accounts, err := repository.GetAccountsByUserID(userID, db)
	if err != nil {
		return 0, err
	}

	// Calculate the total balance of active accounts.
	var totalBalance float64
	for _, account := range accounts {
		if account.IsActive {
			totalBalance += account.Balance
		}
	}

	return totalBalance, nil
}