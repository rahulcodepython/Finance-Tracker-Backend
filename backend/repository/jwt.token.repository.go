package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateJwtToken(db interfaces.SqlExecutor, token *models.JwtToken) error {
	query := `INSERT INTO jwt_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, token.ID, token.UserID, token.Token, token.ExpiresAt, token.CreatedAt)
	return err
}

func GetJwtTokenByUserID(db interfaces.SqlExecutor, userID uuid.UUID) (*models.JwtToken, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM jwt_tokens WHERE user_id = $1`
	row := db.QueryRow(query, userID)

	var token models.JwtToken
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No token found
		}
		return nil, fmt.Errorf("error getting JWT token by user ID: %w", err)
	}

	return &token, nil
}

func GetJwtTokenByToken(db interfaces.SqlExecutor, tokenString string) (*models.JwtToken, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM jwt_tokens WHERE token = $1`
	row := db.QueryRow(query, tokenString)

	var token models.JwtToken
	err := row.Scan(&token.ID, &token.UserID, &token.Token, &token.ExpiresAt, &token.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No token found
		}
		return nil, fmt.Errorf("error getting JWT token by user ID: %w", err)
	}

	return &token, nil
}

func DeleteJwtToken(db interfaces.SqlExecutor, token string) error {
	query := `DELETE FROM jwt_tokens WHERE token = $1`
	_, err := db.Exec(query, token)
	return err
}

func DeleteJwtTokenByUserID(db interfaces.SqlExecutor, userID uuid.UUID) error {
	query := `DELETE FROM jwt_tokens WHERE user_id = $1`
	_, err := db.Exec(query, userID)
	return err
}
