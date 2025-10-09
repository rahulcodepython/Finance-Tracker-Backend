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

func GenerateReport(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	// Get aggregate data
	aggregateData, err := GetAggregateData(userID, startDate, endDate, db)
	if err != nil {
		return nil, err
	}

	// Get spending by category
	spendingByCategory, err := repository.GetSpendingByCategory(userID, db)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"summary":            aggregateData,
		"spendingByCategory": spendingByCategory,
	}, nil
}

func ExportTransactions(userID uuid.UUID, writer io.Writer, db *sql.DB) error {
	transactions, err := repository.GetTransactionsByUserID(userID, db)
	if err != nil {
		return err
	}

	// Fetch category and account names for each transaction
	categoryMap := make(map[uuid.UUID]string)

	categories, err := repository.GetCategories(db)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
	}
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{"ID", "Date", "Description", "Amount", "Type", "Category", "Account"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	// Write rows
	for _, transaction := range transactions {
		row := []string{
			transaction.ID.String(),
			transaction.TransactionDate.Format("2006-01-02"),
			transaction.Description,
			fmt.Sprintf("%.2f", transaction.Amount),
			string(transaction.Type),
		}

		if transaction.CategoryID.Valid {
			row = append(row, categoryMap[transaction.CategoryID.UUID])
		} else {
			row = append(row, "")
		}

		account, err := repository.GetAccountByID(transaction.AccountID, db)
		if err != nil {
			log.Printf("Error fetching account %s: %v", transaction.AccountID, err)
		} else {
			row = append(row, account.Name)
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	return nil
}
