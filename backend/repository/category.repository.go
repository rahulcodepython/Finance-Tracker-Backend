package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateCategory(category *models.Category, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO categories (%s) VALUES ($1, $2, $3)", models.CategoryColumns)
	_, err := db.Exec(query, category.ID, category.Name, category.Type)
	return err
}

func GetCategories(db *sql.DB) ([]models.Category, error) {
	query := "SELECT * FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Type); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetCategoryByID(id uuid.UUID, db *sql.DB) (*models.Category, error) {
	query := "SELECT * FROM categories WHERE id = $1"
	row := db.QueryRow(query, id)

	var category models.Category
	if err := row.Scan(&category.ID, &category.Name, &category.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &category, nil
}

func UpdateCategory(category *models.Category, db *sql.DB) error {
	query := "UPDATE categories SET name = $1, type = $2 WHERE id = $3"
	_, err := db.Exec(query, category.Name, category.Type, category.ID)
	return err
}

func DeleteCategory(id uuid.UUID, db *sql.DB) error {
	query := "DELETE FROM categories WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
