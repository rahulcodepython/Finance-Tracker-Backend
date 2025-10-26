// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateUser inserts a new user record into the database.
func CreateUser(user *models.User, db interfaces.SqlExecutor) error {
	// Construct the SQL query for insertion.
	query := fmt.Sprintf("INSERT INTO users (%s) VALUES ($1, $2, $3, $4, $5, $6)", models.UserColumns)
	// Execute the query with the user data.
	_, err := db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.Provider, user.CreatedAt)
	return err
}

// GetUserByEmail retrieves a single user by their email address from the database.
func GetUserByEmail(email string, db interfaces.SqlExecutor) (*models.User, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM users WHERE email = $1"
	// Execute the query.
	row := db.QueryRow(query, email)
	var user models.User

	// Scan the row into a User struct.
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.CreatedAt); err != nil {
		// If no rows were found, return the error.
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID retrieves a single user by their ID from the database.
func GetUserByID(id uuid.UUID, db interfaces.SqlExecutor) (*models.User, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM users WHERE id = $1"
	// Execute the query.
	row := db.QueryRow(query, id)

	// Scan the row into a User struct.
	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.CreatedAt); err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user record in the database.
func UpdateUser(user *models.User, db interfaces.SqlExecutor) error {
	// Construct the SQL query for update.
	query := "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4"
	// Execute the query with the updated user data.
	_, err := db.Exec(query, user.Name, user.Email, user.Password, user.ID)
	return err
}