// Package scheduler provides functionality for scheduling and running background tasks.
package scheduler

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-co-op/gocron" // Go-cron library for scheduling tasks
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// StartScheduler initializes and starts the cron scheduler for recurring tasks.
func StartScheduler(db *sql.DB) {
	// Create a new scheduler that uses the local timezone.
	s := gocron.NewScheduler(time.Local)

	// Schedule the ProcessRecurringTransactions function to run every day at midnight.
	s.Every(1).Day().At("00:00").Do(func() {
		log.Println("Running recurring transaction check...")
		ProcessRecurringTransactions(db)
		log.Println("Complete for today.")
	})

	// Start the scheduler asynchronously.
	s.StartAsync()
}

// ProcessRecurringTransactions retrieves and processes all recurring transactions.
func ProcessRecurringTransactions(db *sql.DB) {
	// Get all recurring transactions from the database.
	recurringTransactions, err := repository.GetRecurringTransactions(db)
	if err != nil {
		log.Println("Error getting recurring transactions:", err)
		return
	}

	// Iterate over each recurring transaction and create a new transaction if it is due today.
	for _, rt := range recurringTransactions {
		today := time.Now().In(utils.LOC).Day()

		if rt.RecurringFrequency == models.Monthly && rt.RecurringDate == today {
			createTransactionFromRecurring(rt, db)
		} else if rt.RecurringFrequency == models.Yearly && rt.RecurringDate == today && time.Now().In(utils.LOC).Month() == rt.CreatedAt.Month() {
			createTransactionFromRecurring(rt, db)
		}
	}
}

// createTransactionFromRecurring creates a new transaction from a recurring transaction.
func createTransactionFromRecurring(rt models.RecurringTransaction, db *sql.DB) {
	// Create a new transaction model from the recurring transaction data.
	transaction := &models.Transaction{
		ID:              uuid.New(),
		UserID:          rt.UserID,
		AccountID:       rt.AccountID,
		CategoryID:      rt.CategoryID,
		BudgetID:        rt.BudgetID,
		Description:     rt.Description,
		Amount:          rt.Amount,
		Type:            rt.Type,
		Note:            rt.Note,
		TransactionDate: time.Now().In(utils.LOC),
		CreatedAt:       time.Now().In(utils.LOC),
		UpdatedAt:       time.Now().In(utils.LOC),
	}

	// Create the new transaction in the database.
	if err := repository.CreateTransaction(transaction, db); err != nil {
		log.Println("Error creating transaction from recurring:", err)
	}
}