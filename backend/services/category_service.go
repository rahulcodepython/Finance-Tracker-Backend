package services

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func CreateCategory(name string, categoryType models.TransactionType, db *sql.DB) (*models.Category, error) {
	category := &models.Category{
		ID:   uuid.New(),
		Name: name,
		Type: categoryType,
	}

	err := repository.CreateCategory(category, db)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func GetCategories(db *sql.DB) ([]models.Category, error) {
	return repository.GetCategories(db)
}

func UpdateCategory(id uuid.UUID, name string, categoryType models.TransactionType, db *sql.DB) (*models.Category, error) {
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return nil, err
	}

	category.Name = name
	category.Type = categoryType

	err = repository.UpdateCategory(category, db)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func DeleteCategory(id uuid.UUID, db *sql.DB) error {
	return repository.DeleteCategory(id, db)
}
