package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
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

	// Log the creation
	go CreateLog(userID, fmt.Sprintf("New budget '%s' created", budget.Name), db)

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

	if budget == nil {
		return nil, sql.ErrNoRows
	}

	budget.Name = name
	budget.Amount = amount
	budget.UpdatedAt = time.Now().In(utils.LOC)

	err = repository.UpdateBudget(budget, db)
	if err != nil {
		return nil, err
	}

	// Log the update
	go CreateLog(budget.UserID, fmt.Sprintf("Budget '%s' updated", budget.Name), db)

	return budget, nil
}

func DeleteBudget(id uuid.UUID, db *sql.DB) error {
	budget, err := repository.GetBudgetByID(id, db)
	if err != nil {
		return err
	}

	if budget == nil {
		return sql.ErrNoRows
	}

	err = repository.DeleteBudget(id, db)
	if err != nil {
		return err
	}

	// Log the deletion
	go CreateLog(budget.UserID, fmt.Sprintf("Budget '%s' removed", budget.Name), db)

	return nil
}

func CheckBudgetExistsById(id uuid.UUID, db *sql.DB) (bool, error) {
	budget, err := repository.GetBudgetByID(id, db)
	if err != nil {
		return false, err
	}

	if budget != nil {
		return true, nil
	}

	return false, nil
}
