// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateCategory inserts a new category record into the database.
func CreateCategory(category *models.Category, db interfaces.SqlExecutor) error {
	// Construct the SQL query for insertion.
	query := fmt.Sprintf("INSERT INTO categories (%s) VALUES ($1, $2, $3)", models.CategoryColumns)
	// Execute the query with the category data.
	_, err := db.Exec(query, category.ID, category.Name, category.Type)
	return err
}

// GetCategories retrieves all categories from the database.
func GetCategories(db interfaces.SqlExecutor) ([]models.Category, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM categories"
	// Execute the query.
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into Category structs.
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

// GetCategoryByID retrieves a single category by its ID from the database.
func GetCategoryByID(id uuid.UUID, db interfaces.SqlExecutor) (*models.Category, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM categories WHERE id = $1"
	// Execute the query.
	row := db.QueryRow(query, id)

	// Scan the row into a Category struct.
	var category models.Category
	if err := row.Scan(&category.ID, &category.Name, &category.Type); err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// UpdateCategory updates an existing category record in the database.
func UpdateCategory(category *models.Category, db interfaces.SqlExecutor) error {
	// Construct the SQL query for update.
	query := "UPDATE categories SET name = $1, type = $2 WHERE id = $3"
	// Execute the query with the updated category data.
	_, err := db.Exec(query, category.Name, category.Type, category.ID)
	return err
}

// DeleteCategory deletes a category record from the database by its ID.
func DeleteCategory(id uuid.UUID, db interfaces.SqlExecutor) error {
	// Construct the SQL query for deletion.
	query := "DELETE FROM categories WHERE id = $1"
	// Execute the query.
	_, err := db.Exec(query, id)
	return err
}