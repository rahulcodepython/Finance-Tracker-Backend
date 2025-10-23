package services

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func GetDashboardSummary(userID uuid.UUID, db *sql.DB) (map[string]interface{}, error) {
	// Get total balance
	totalBalance, err := GetTotalBalance(userID, db)
	if err != nil {
		return nil, err
	}

	// Get monthly income and expenses
	// This is a simplified example. In a real application, you would have a more complex query to get the data for the current month.
	aggregateData, err := GetAggregateData(userID, "", "", db)
	if err != nil {
		return nil, err
	}

	// Get recent transactions
	recentTransactions, err := GetTransactions(userID, 1, 10, "", "", "", "", "", "", db)
	if err != nil {
		return nil, err
	}

	// Get monthly spending by category
	// This is a simplified example. In a real application, you would have a more complex query to get this data.
	spendingByCategory, err := repository.GetSpendingByCategory(userID, db)
	if err != nil {
		return nil, err
	}

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
		},
		"recentTransactions": recentTransactions,
	}, nil
}
