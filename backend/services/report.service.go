// Package services provides business logic for the application.
package services

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

// GenerateReport generates a financial report for a user within a specified date range.
func GenerateReport(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	// Get aggregate data for the specified date range.
	aggregateData, err := GetAggregateData(userID, startDate, endDate, db)
	if err != nil {
		return nil, err
	}

	// Get spending by category.
	spendingByCategory, err := repository.GetSpendingByCategory(userID, db)
	if err != nil {
		return nil, err
	}

	// Return the report data in a map.
	return map[string]interface{}{
		"summary":            aggregateData,
		"spendingByCategory": spendingByCategory,
	}, nil
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
	categories, err := repository.GetCategories(db)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
	}
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
		account, err := repository.GetAccountByID(transaction.AccountID, db)
		if err != nil {
			log.Printf("Error fetching account %s: %v", transaction.AccountID, err)
		}

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