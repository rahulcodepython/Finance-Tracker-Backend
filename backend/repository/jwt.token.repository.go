// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateJwtToken inserts a new JWT token record into the database.
func CreateJwtToken(db interfaces.SqlExecutor, token *models.JwtToken) error {
	// Construct the SQL query for insertion.
	query := `INSERT INTO jwt_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	// Execute the query with the token data.
	_, err := db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	return err
}

// GetJwtTokenByUserID retrieves a JWT token by user ID from the database.
func GetJwtTokenByUserID(db interfaces.SqlExecutor, userID uuid.UUID) (*models.JwtToken, error) {
	// Construct the SQL query for selection.
	query := `SELECT id, user_id, token, expires_at, created_at FROM jwt_tokens WHERE user_id = $1`
	// Execute the query.
	row := db.QueryRow(query, userID)

	// Scan the row into a JwtToken struct.
	var token models.JwtToken
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting JWT token by user ID: %w", err)
	}

	return &token, nil
}

// GetJwtTokenByToken retrieves a JWT token by the token string from the database.
func GetJwtTokenByToken(db interfaces.SqlExecutor, tokenString string) (*models.JwtToken, error) {
	// Construct the SQL query for selection.
	query := `SELECT id, user_id, token, expires_at, created_at FROM jwt_tokens WHERE token = $1`
	// Execute the query.
	row := db.QueryRow(query, tokenString)

	// Scan the row into a JwtToken struct.
	var token models.JwtToken
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting JWT token by token string: %w", err)
	}

	return &token, nil
}

// DeleteJwtToken deletes a JWT token from the database by the token string.
func DeleteJwtToken(db interfaces.SqlExecutor, token string) error {
	// Construct the SQL query for deletion.
	query := `DELETE FROM jwt_tokens WHERE token = $1`
	// Execute the query.
	_, err := db.Exec(query, token)
	return err
}

// DeleteJwtTokenByUserID deletes all JWT tokens for a given user ID from the database.
func DeleteJwtTokenByUserID(db interfaces.SqlExecutor, userID uuid.UUID) error {
	// Construct the SQL query for deletion.
	query := `DELETE FROM jwt_tokens WHERE user_id = $1`
	// Execute the query.
	_, err := db.Exec(query, userID)
	return err
}