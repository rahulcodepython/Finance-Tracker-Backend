package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateBudget(budget *models.Budget, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO budgets (%s) VALUES ($1, $2, $3, $4, $5, $6)", models.BudgetColumns)
	_, err := db.Exec(query, budget.ID, budget.UserID, budget.Name, budget.Amount, budget.CreatedAt, budget.UpdatedAt)
	return err
}

func GetBudgetsByUserID(userID uuid.UUID, db *sql.DB) ([]models.Budget, error) {
	query := "SELECT * FROM budgets WHERE user_id = $1"
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func GetBudgetByID(id uuid.UUID, db *sql.DB) (*models.Budget, error) {
	query := "SELECT * FROM budgets WHERE id = $1"
	row := db.QueryRow(query, id)

	var budget models.Budget
	if err := row.Scan(&budget.ID, &budget.UserID, &budget.Name, &budget.Amount, &budget.CreatedAt, &budget.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}
	return &budget, nil
}

func UpdateBudget(budget *models.Budget, db *sql.DB) error {
	query := "UPDATE budgets SET name = $1, amount = $2, updated_at = $3 WHERE id = $4"
	_, err := db.Exec(query, budget.Name, budget.Amount, budget.UpdatedAt, budget.ID)
	return err
}

func DeleteBudget(id uuid.UUID, db *sql.DB) error {
	query := "DELETE FROM budgets WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
