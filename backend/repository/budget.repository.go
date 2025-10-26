// Package repository provides functionality for interacting with the database.
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateBudget inserts a new budget record into the database.
func CreateBudget(budget *models.Budget, db interfaces.SqlExecutor) error {
	// Construct the SQL query for insertion.
	query := fmt.Sprintf("INSERT INTO budgets (%s) VALUES ($1, $2, $3, $4, $5, $6)", models.BudgetColumns)
	// Execute the query with the budget data.
	_, err := db.Exec(query, budget.ID, budget.UserID, budget.Name, budget.Amount, budget.CreatedAt, budget.UpdatedAt)
	return err
}

// GetBudgetsByUserID retrieves all budgets for a given user ID from the database.
func GetBudgetsByUserID(userID uuid.UUID, db interfaces.SqlExecutor) ([]models.Budget, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM budgets WHERE user_id = $1"
	// Execute the query.
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into Budget structs.
	var budgets []models.Budget
	for rows.Next() {
		var budget models.Budget
		if err := rows.Scan(&budget.ID, &budget.UserID, &budget.Name, &budget.Amount, &budget.CreatedAt, &budget.UpdatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}

// GetBudgetByID retrieves a single budget by its ID from the database.
func GetBudgetByID(id uuid.UUID, db interfaces.SqlExecutor) (*models.Budget, error) {
	// Construct the SQL query for selection.
	query := "SELECT * FROM budgets WHERE id = $1"
	// Execute the query.
	row := db.QueryRow(query, id)

	// Scan the row into a Budget struct.
	var budget models.Budget
	if err := row.Scan(&budget.ID, &budget.UserID, &budget.Name, &budget.Amount, &budget.CreatedAt, &budget.UpdatedAt); err != nil {
		// If no rows were found, return nil.
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &budget, nil
}

// UpdateBudget updates an existing budget record in the database.
func UpdateBudget(budget *models.Budget, db interfaces.SqlExecutor) error {
	// Construct the SQL query for update.
	query := "UPDATE budgets SET name = $1, amount = $2, updated_at = $3 WHERE id = $4"
	// Execute the query with the updated budget data.
	_, err := db.Exec(query, budget.Name, budget.Amount, budget.UpdatedAt, budget.ID)
	return err
}

// DeleteBudget deletes a budget record from the database by its ID.
func DeleteBudget(id uuid.UUID, db interfaces.SqlExecutor) error {
	// Construct the SQL query for deletion.
	query := "DELETE FROM budgets WHERE id = $1"
	// Execute the query.
	_, err := db.Exec(query, id)
	return err
}