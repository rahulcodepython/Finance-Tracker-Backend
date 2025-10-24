package services

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func GetDashboardSummary(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, budgetID string, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	// Get total balance
	totalBalance, err := GetTotalBalance(userID, db)
	if err != nil {
		return nil, err
	}

	aggregateData, err := GetAggregateData(userID, startDate, endDate, db)
	if err != nil {
		return nil, err
	}

	// Get recent transactions
	recentTransactions, err := GetTransactions(userID, page, limit, description, categoryID, accountID, budgetID, startDate, endDate, db)
	if err != nil {
		return nil, err
	}

	spendingByCategory, err := repository.GetSpendingByCategory(userID, db)
	if err != nil {
		return nil, err
	}

	earningByCategory, err := repository.GetEarningByCategory(userID, db)
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
			"earningByCategory":  earningByCategory,
		},
		"recentTransactions": recentTransactions,
	}, nil
}
