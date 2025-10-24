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

	if category == nil {
		return nil, sql.ErrNoRows
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
	category, err := repository.GetCategoryByID(id, db)
	if err != nil {
		return err
	}

	if category == nil {
		return sql.ErrNoRows
	}

	return repository.DeleteCategory(id, db)
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
