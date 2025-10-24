package services

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func CreateCategory(name string, categoryType models.TransactionType, userID uuid.UUID, db *sql.DB) (*models.Category, error) {
	category := &models.Category{
		ID:   uuid.New(),
		Name: name,
		Type: categoryType,
	}

	err := repository.CreateCategory(category, db)
	if err != nil {
		return nil, err
	}

	// Log the creation
	go CreateLog(userID, fmt.Sprintf("New category '%s' created", category.Name), db)

	return category, nil
}

func GetCategories(db *sql.DB) ([]models.Category, error) {
	return repository.GetCategories(db)
}

func UpdateCategory(id uuid.UUID, name string, categoryType models.TransactionType, userID uuid.UUID, db *sql.DB) (*models.Category, error) {
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, sql.ErrNoRows
	}

	category.Name = name
	category.Type = categoryType

	err = repository.UpdateCategory(category, db)
	if err != nil {
		return nil, err
	}

	// Log the update
	go CreateLog(userID, fmt.Sprintf("Category '%s' updated", category.Name), db)

	return category, nil
}

func DeleteCategory(id uuid.UUID, userID uuid.UUID, db *sql.DB) error {
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return err
	}

	if category == nil {
		return sql.ErrNoRows
	}

	err = repository.DeleteCategory(id, db)
	if err != nil {
		return err
	}

	// Log the deletion
	go CreateLog(userID, fmt.Sprintf("Category '%s' removed", category.Name), db)

	return nil
}

func CheckCategoryExistsById(id uuid.UUID, db *sql.DB) (bool, error) {
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return false, err
	}

	if category != nil {
		return true, nil
	}

	return false, nil
}
