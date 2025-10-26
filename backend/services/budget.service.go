// Package services provides business logic for the application.
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

// CreateBudget creates a new budget for a user.
func CreateBudget(userID uuid.UUID, name string, amount float64, db *sql.DB) (*models.Budget, error) {
	// Create a new Budget model.
	budget := &models.Budget{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Amount:    amount,
		CreatedAt: time.Now().In(utils.LOC),
		UpdatedAt: time.Now().In(utils.LOC),
	}

	// Create the budget in the database.
	err := repository.CreateBudget(budget, db)
	if err != nil {
		return nil, err
	}

	// Create a log entry for the budget creation.
	go CreateLog(userID, fmt.Sprintf("New budget '%s' created", budget.Name), db)

	return budget, nil
}

// GetBudgets retrieves all budgets for a user.
func GetBudgets(userID uuid.UUID, db *sql.DB) ([]models.Budget, error) {
	return repository.GetBudgetsByUserID(userID, db)
}

// UpdateBudget updates an existing budget.
func UpdateBudget(id uuid.UUID, name string, amount float64, db *sql.DB) (*models.Budget, error) {
	// Get the budget from the database.
	budget, err := repository.GetBudgetByID(id, db)
	if err != nil {
		return nil, err
	}

	// If the budget does not exist, return an error.
	if budget == nil {
		return nil, sql.ErrNoRows
	}

	// Update the budget fields.
	budget.Name = name
	budget.Amount = amount
	budget.UpdatedAt = time.Now().In(utils.LOC)

	// Update the budget in the database.
	err = repository.UpdateBudget(budget, db)
	if err != nil {
		return nil, err
	}

	// Create a log entry for the budget update.
	go CreateLog(budget.UserID, fmt.Sprintf("Budget '%s' updated", budget.Name), db)

	return budget, nil
}

// DeleteBudget deletes a budget.
func DeleteBudget(id uuid.UUID, db *sql.DB) error {
	// Get the budget from the database.
	budget, err := repository.GetBudgetByID(id, db)
	if err != nil {
		return err
	}

	// If the budget does not exist, return an error.
	if budget == nil {
		return sql.ErrNoRows
	}

	// Delete the budget from the database.
	err = repository.DeleteBudget(id, db)
	if err != nil {
		return err
	}

	// Create a log entry for the budget deletion.
	go CreateLog(budget.UserID, fmt.Sprintf("Budget '%s' removed", budget.Name), db)

	return nil
}