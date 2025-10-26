// Package services provides business logic for the application.
package services

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

// GetDashboardSummary retrieves a summary of the user's financial data for the dashboard.
func GetDashboardSummary(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, budgetID string, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	// Get the total balance of all active accounts.
	totalBalance, err := GetTotalBalance(userID, db)
	if err != nil {
		return nil, err
	}

	// Get aggregate data for the specified date range.
	aggregateData, err := GetAggregateData(userID, startDate, endDate, db)
	if err != nil {
		return nil, err
	}

	// Get recent transactions with optional filters and pagination.
	recentTransactions, err := GetTransactions(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
	if err != nil {
		return nil, err
	}

	// Get spending by category.
	spendingByCategory, err := repository.GetSpendingByCategory(userID, db)
	if err != nil {
		return nil, err
	}

	// Get earnings by category.
	earningByCategory, err := repository.GetEarningByCategory(userID, db)
	if err != nil {
		return nil, err
	}

	// Return the dashboard summary in a map.
	return map[string]interface{}{
		"summary": map[string]interface{}{
			"totalBalance":    totalBalance,
			"monthlyIncome":   aggregateData["totalIncome"],
			"monthlyExpenses": aggregateData["totalExpenses"],
			"monthlySavings":  aggregateData["netIncome"],
		},
		"graphs": map[string]interface{}{
			"incomeVsExpense":    []map[string]interface{}{},
			"spendingByCategory": spendingByCategory,
			"earningByCategory":  earningByCategory,
		},
		"recentTransactions": recentTransactions,
	}, nil
}