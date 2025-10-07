package scheduler

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func Start(db *sql.DB) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Day().At("00:00").Do(func() {
		log.Println("Running recurring transactions job")
		processRecurringTransactions(db)
	})

	s.StartAsync()
}

func processRecurringTransactions(db *sql.DB) {
	rows, err := db.Query("SELECT id, user_id, account_id, category_id, description, amount, type, frequency, start_date, end_date FROM recurring_transactions WHERE is_active = true AND start_date <= NOW() AND (end_date IS NULL OR end_date >= NOW())")
	if err != nil {
		log.Printf("Error fetching recurring transactions: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rt models.RecurringTransaction
		if err := rows.Scan(&rt.ID, &rt.UserID, &rt.AccountID, &rt.CategoryID, &rt.Description, &rt.Amount, &rt.Type, &rt.Frequency, &rt.StartDate, &rt.EndDate); err != nil {
			log.Printf("Error scanning recurring transaction: %v", err)
			continue
		}

		// Check if a transaction has already been created for the current period
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM transactions WHERE recurring_transaction_id = $1 AND transaction_date = $2", rt.ID, time.Now().Format("2006-01-02")).Scan(&count)
		if err != nil {
			log.Printf("Error checking for existing transaction: %v", err)
			continue
		}

		if count == 0 {
			transaction := &models.Transaction{
				ID:                    uuid.New(),
				UserID:                rt.UserID,
				AccountID:             rt.AccountID,
				CategoryID:            rt.CategoryID,
				Description:           rt.Description,
				Amount:                rt.Amount,
				Type:                  rt.Type,
				TransactionDate:       time.Now(),
				RecurringTransactionID: uuid.NullUUID{UUID: rt.ID, Valid: true},
				CreatedAt:             time.Now(),
				UpdatedAt:             time.Now(),
			}

			if err := repository.CreateTransaction(transaction, db); err != nil {
				log.Printf("Error creating transaction: %v", err)
				continue
			}

			// Update account balance
			account, err := repository.GetAccountByID(rt.AccountID, db)
			if err != nil {
				log.Printf("Error getting account: %v", err)
				continue
			}

			if rt.Type == models.TransactionTypeIncome {
				account.Balance += rt.Amount
			} else {
				account.Balance -= rt.Amount
			}

			if err := repository.UpdateAccount(account, db); err != nil {
				log.Printf("Error updating account: %v", err)
				continue
			}
		}
	}
}
