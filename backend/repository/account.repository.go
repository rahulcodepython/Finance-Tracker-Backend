package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateAccount(account *models.Account, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO accounts (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", models.AccountColumns)
	_, err := db.Exec(query, account.ID, account.UserID, account.Name, account.Type, account.Balance, account.IsActive, account.CreatedAt, account.UpdatedAt)
	return err
}

func GetAccountsByUserID(userID uuid.UUID, db *sql.DB) ([]models.Account, error) {
	query := "SELECT * FROM accounts WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func GetAccountByID(id uuid.UUID, db *sql.DB) (*models.Account, error) {
	query := "SELECT * FROM accounts WHERE id = $1"
	row := db.QueryRow(query, id)

	var account models.Account
	if err := row.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Balance, &account.IsActive, &account.CreatedAt, &account.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &account, nil
}

func UpdateAccount(account *models.Account, db *sql.DB) error {
	query := "UPDATE accounts SET name = $1, type = $2, balance = $3, is_active = $4, updated_at = $5 WHERE id = $6"
	_, err := db.Exec(query, account.Name, account.Type, account.Balance, account.IsActive, account.UpdatedAt, account.ID)
	return err
}

func DeleteAccount(id uuid.UUID, db *sql.DB) error {
	query := "DELETE FROM accounts WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
