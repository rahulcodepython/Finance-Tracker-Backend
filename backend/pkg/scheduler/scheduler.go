package scheduler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func StartScheduler(db *sql.DB) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Day().At("00:00").Do(func() {
		fmt.Println("Running recurring transaction check...")
		ProcessRecurringTransactions(db)
	})

	s.StartAsync()
}

func ProcessRecurringTransactions(db *sql.DB) {
	query := `SELECT id, user_id, account_id, category_id, description, amount, type, recurring_frequency, recurring_date FROM recurring_transactions`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error getting recurring transactions:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rt models.RecurringTransaction
		if err := rows.Scan(&rt.ID, &rt.UserID, &rt.AccountID, &rt.CategoryID, &rt.Description, &rt.Amount, &rt.Type, &rt.RecurringFrequency, &rt.RecurringDate); err != nil {
			fmt.Println("Error scanning recurring transaction:", err)
			continue
		}

		today := time.Now().Day()

		if rt.RecurringFrequency == models.Monthly && rt.RecurringDate == today {
			createTransactionFromRecurring(rt, db)
		} else if rt.RecurringFrequency == models.Yearly && rt.RecurringDate == today && time.Now().Month() == rt.CreatedAt.Month() {
			createTransactionFromRecurring(rt, db)
		}
	}
}

func createTransactionFromRecurring(rt models.RecurringTransaction, db *sql.DB) {
	t := &models.Transaction{
		ID:              uuid.New(),
		UserID:          rt.UserID,
		AccountID:       rt.AccountID,
		CategoryID:      rt.CategoryID,
		Description:     rt.Description,
		Amount:          rt.Amount,
		Type:            rt.Type,
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := repository.CreateTransaction(t, db); err != nil {
		fmt.Println("Error creating transaction from recurring:", err)
	}
}
