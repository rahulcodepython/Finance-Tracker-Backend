package services

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func CreateBudget(userID uuid.UUID, name string, amount float64, db *sql.DB) (*models.Budget, error) {
	budget := &models.Budget{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
		Amount: amount,
	}

	err := repository.CreateBudget(budget, db)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func GetBudgets(userID uuid.UUID, db *sql.DB) ([]models.Budget, error) {
	return repository.GetBudgetsByUserID(userID, db)
}

func UpdateBudget(id uuid.UUID, name string, amount float64, db *sql.DB) (*models.Budget, error) {
	budget, err := repository.GetBudgetByID(id, db)
	if err != nil {
		return nil, err
	}

	budget.Name = name
	budget.Amount = amount

	err = repository.UpdateBudget(budget, db)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func DeleteBudget(id uuid.UUID, db *sql.DB) error {
	return repository.DeleteBudget(id, db)
}