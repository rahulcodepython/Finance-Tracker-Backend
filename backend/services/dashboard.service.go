// Package services provides business logic for the application.
package services

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/serializers"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// GetDashboardSummary retrieves a summary of the user's financial data for the dashboard.
func GetDashboardSummary(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (*serializers.DashboardResponse, error) {
	if startDate == "" {
		startDate = time.Now().In(utils.LOC).AddDate(0, 0, -30).Format("2006-01-02")
	}

	if endDate == "" {
		endDate = time.Now().In(utils.LOC).Format("2006-01-02")
	}

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
	aggregateData.TotalBalance = totalBalance

	// Get recent transactions with optional filters and pagination.
	recentTransactions, err := GetTransactions(userID, 1, 10, "", "", "", "", startDate, endDate, db)
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

	getIncomeExpense, err := repository.GetIncomeExpense(userID, db)
	if err != nil {
		return nil, err
	}

	var dashboardResponse serializers.DashboardResponse
	dashboardResponse.Summary = *aggregateData
	dashboardResponse.Graphs.IncomeExpenseAggregate = *getIncomeExpense
	dashboardResponse.Graphs.SpendingByCategory = *spendingByCategory
	dashboardResponse.Graphs.EarningByCategory = *earningByCategory
	dashboardResponse.RecentTransactions = recentTransactions

	return &dashboardResponse, nil
}

// ExportTransactions exports a user's transactions to a CSV file.
func ExportTransactions(userID uuid.UUID, writer io.Writer, db *sql.DB) error {
	// Get all transactions for the user.
	transactions, err := repository.GetTransactionsByUserID(userID, db)
	if err != nil {
		return err
	}

	// Create a map of category IDs to category names for efficient lookup.
	categoryMap := make(map[uuid.UUID]string)
	categories, _ := repository.GetCategories(db)

	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	// Create a new CSV writer.
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write the CSV header.
	header := []string{"ID", "Date", "Description", "Category", "Amount", "Type", "Account"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	// Write each transaction as a row in the CSV file.
	for _, transaction := range transactions {
		// Get the account name for the transaction.
		account, _ := repository.GetAccountByID(transaction.AccountID, db)

		// Construct the row with transaction data.
		row := []string{
			transaction.ID.String(),
			transaction.TransactionDate.Format("2006-01-02"),
			transaction.Description,
			categoryMap[transaction.CategoryID],
			fmt.Sprintf("%.2f", transaction.Amount),
			string(transaction.Type),
			account.Name,
		}

		// Write the row to the CSV file.
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}
