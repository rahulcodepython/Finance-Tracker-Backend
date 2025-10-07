package repository

import (
	"database/sql"
	"fmt"

	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateUser(user *models.User, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO users (%s) VALUES ($1, $2, $3, $4, $5, $6)", models.UserColumns)
	_, err := db.Exec(query, user.ID, user.Name, user.Email, user.Password, user.Provider, user.CreatedAt)
	return err
}

func GetUserByEmail(email string, db *sql.DB) (*models.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	row := db.QueryRow(query, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id string, db *sql.DB) (*models.User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	row := db.QueryRow(query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Provider, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *models.User, db *sql.DB) error {
	query := "UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4"
	_, err := db.Exec(query, user.Name, user.Email, user.Password, user.ID)
	return err
}
