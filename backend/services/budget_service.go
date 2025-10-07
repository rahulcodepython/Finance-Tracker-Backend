package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateBudget(userID uuid.UUID, categoryID uuid.UUID, amount float64, month time.Time, db *sql.DB) (*models.Budget, error) {
	return nil, nil
}

func GetBudgets(userID uuid.UUID, db *sql.DB) ([]models.Budget, error) {
	return nil, nil
}

func UpdateBudget(id uuid.UUID, categoryID uuid.UUID, amount float64, month time.Time, db *sql.DB) (*models.Budget, error) {
	return nil, nil
}

func DeleteBudget(id uuid.UUID, db *sql.DB) error {
	return nil
}
