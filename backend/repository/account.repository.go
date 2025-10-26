// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateAccount inserts a new account record into the database.
func CreateAccount(account *models.Account, db interfaces.SqlExecutor) error {
	// Construct the SQL query for insertion.
	query := fmt.Sprintf("INSERT INTO accounts (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", models.AccountColumns)
	// Execute the query with the account data.
	_, err := db.Exec(query, account.ID, account.UserID, account.Name, account.Type, account.Balance, account.IsActive, account.CreatedAt, account.UpdatedAt)
	return err
}

// GetAccountsByUserID retrieves all accounts for a given user ID from the database.
func GetAccountsByUserID(userID uuid.UUID, db interfaces.SqlExecutor) ([]models.Account, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM accounts WHERE user_id = $1"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into Account structs.
	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Balance, &account.IsActive, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// GetAccountByID retrieves a single account by its ID from the database.
func GetAccountByID(id uuid.UUID, db interfaces.SqlExecutor) (*models.Account, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM accounts WHERE id = $1"
	// Execute the query.
	row := db.QueryRow(query, id)

	// Scan the row into an Account struct.
	var account models.Account
	if err := row.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Balance, &account.IsActive, &account.CreatedAt, &account.UpdatedAt); err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// UpdateAccount updates an existing account record in the database.
func UpdateAccount(account *models.Account, db interfaces.SqlExecutor) error {
	// Construct the SQL query for update.
	query := "UPDATE accounts SET name = $1, type = $2, balance = $3, is_active = $4, updated_at = $5 WHERE id = $6"
	// Execute the query with the updated account data.
	_, err := db.Exec(query, account.Name, account.Type, account.Balance, account.IsActive, account.UpdatedAt, account.ID)
	return err
}

// DeleteAccount deletes an account record from the database by its ID.
func DeleteAccount(id uuid.UUID, db interfaces.SqlExecutor) error {
	// Construct the SQL query for deletion.
	query := "DELETE FROM accounts WHERE id = $1"
	// Execute the query.
	_, err := db.Exec(query, id)
	return err
}