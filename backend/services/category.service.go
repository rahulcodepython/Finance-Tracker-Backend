// Package services provides business logic for the application.
package services

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

// CreateCategory creates a new transaction category.
func CreateCategory(name string, categoryType models.TransactionType, userID uuid.UUID, db *sql.DB) (*models.Category, error) {
	// Create a new Category model.
	category := &models.Category{
		ID:   uuid.New(),
		Name: name,
		Type: categoryType,
	}

	// Create the category in the database.
	err := repository.CreateCategory(category, db)
	if err != nil {
		return nil, err
	}

	// Create a log entry for the category creation.
	go CreateLog(userID, fmt.Sprintf("New category '%s' created", category.Name), db)

	return category, nil
}

// GetCategories retrieves all transaction categories.
func GetCategories(db *sql.DB) ([]models.Category, error) {
	return repository.GetCategories(db)
}

// UpdateCategory updates an existing transaction category.
func UpdateCategory(id uuid.UUID, name string, categoryType models.TransactionType, userID uuid.UUID, db *sql.DB) (*models.Category, error) {
	// Get the category from the database.
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return nil, err
	}

	// If the category does not exist, return an error.
	if category == nil {
		return nil, sql.ErrNoRows
	}

	// Update the category fields.
	category.Name = name
	category.Type = categoryType

	// Update the category in the database.
	err = repository.UpdateCategory(category, db)
	if err != nil {
		return nil, err
	}

	// Create a log entry for the category update.
	go CreateLog(userID, fmt.Sprintf("Category '%s' updated", category.Name), db)

	return category, nil
}

// DeleteCategory deletes a transaction category.
func DeleteCategory(id uuid.UUID, userID uuid.UUID, db *sql.DB) error {
	// Get the category from the database.
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return err
	}

	// If the category does not exist, return an error.
	if category == nil {
		return sql.ErrNoRows
	}

	// Delete the category from the database.
	err = repository.DeleteCategory(id, db)
	if err != nil {
		return err
	}

	// Create a log entry for the category deletion.
	go CreateLog(userID, fmt.Sprintf("Category '%s' removed", category.Name), db)

	return nil
}

// CheckCategoryExistsById checks if a category exists by its ID.
func CheckCategoryExistsById(id uuid.UUID, db *sql.DB) (bool, error) {
	// Get the category from the database.
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return false, err
	}

	// If the category is not nil, it exists.
	if category != nil {
		return true, nil
	}

	return false, nil
}