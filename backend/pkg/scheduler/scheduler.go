package scheduler

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func StartScheduler(db *sql.DB) {
	s := gocron.NewScheduler(time.Local)

	s.Every(1).Day().At("00:00").Do(func() {
		log.Println("Running recurring transaction check...")
		ProcessRecurringTransactions(db)
		log.Println("Complete for today.")
	})

	s.StartAsync()
}

func ProcessRecurringTransactions(db *sql.DB) {
	recurringTransactions, err := repository.GetRecurringTransactions(db)
	if err != nil {
		log.Println("Error getting recurring transactions:", err)
		return
	}

	for _, rt := range recurringTransactions {
		today := time.Now().In(utils.LOC).Day()

		if rt.RecurringFrequency == models.Monthly && rt.RecurringDate == today {
			createTransactionFromRecurring(rt, db)
		} else if rt.RecurringFrequency == models.Yearly && rt.RecurringDate == today && time.Now().In(utils.LOC).Month() == rt.CreatedAt.Month() {
			createTransactionFromRecurring(rt, db)
		}
	}
}

func createTransactionFromRecurring(rt models.RecurringTransaction, db *sql.DB) {
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

	if err := repository.CreateTransaction(transaction, db); err != nil {
		log.Println("Error creating transaction from recurring:", err)
	}
}
