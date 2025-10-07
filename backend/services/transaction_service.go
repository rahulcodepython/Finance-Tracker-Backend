package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
)

func CreateTransaction(userID uuid.UUID, accountID uuid.UUID, categoryID uuid.NullUUID, description string, amount float64, transactionType models.TransactionType, transactionDate time.Time, note sql.NullString, db *sql.DB) (*models.Transaction, error) {
	transaction := &models.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		AccountID:       accountID,
		CategoryID:      categoryID,
		Description:     description,
		Amount:          amount,
		Type:            transactionType,
		TransactionDate: transactionDate,
		Note:            note,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback on error, commit on success

	// Create the transaction
	if err := repository.CreateTransaction(transaction, db); err != nil {
		return nil, err
	}

	// Update account balance
	account, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}
	if transactionType == models.TransactionTypeIncome {
		account.Balance += amount
	} else {
		account.Balance -= amount
	}
	if err := repository.UpdateAccount(account, db); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func GetTransactions(userID uuid.UUID, page int, limit int, description string, categoryID string, accountID string, startDate string, endDate string, db *sql.DB) ([]models.Transaction, error) {
	return repository.GetTransactionsByUserIDWithFilters(userID, page, limit, description, categoryID, accountID, startDate, endDate, db)
}

func UpdateTransaction(id uuid.UUID, accountID uuid.UUID, categoryID uuid.NullUUID, description string, amount float64, transactionType models.TransactionType, transactionDate time.Time, note sql.NullString, db *sql.DB) (*models.Transaction, error) {
	transaction, err := repository.GetTransactionByID(id, db)
	if err != nil {
		return nil, err
	}

	oldAccountID := transaction.AccountID
	oldAmount := transaction.Amount
	oldType := transaction.Type // Store old type for balance reversal

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback on error, commit on success

	// Revert old transaction amount from old account balance
	oldAccount, err := repository.GetAccountByID(oldAccountID, db)
	if err != nil {
		return nil, err
	}
	if oldType == models.TransactionTypeIncome {
		oldAccount.Balance -= oldAmount
	} else {
		oldAccount.Balance += oldAmount
	}
	if err := repository.UpdateAccount(oldAccount, db); err != nil {
		return nil, err
	}

	// Update transaction details
	transaction.AccountID = accountID
	transaction.CategoryID = categoryID
	transaction.Description = description
	transaction.Amount = amount
	transaction.Type = transactionType
	transaction.TransactionDate = transactionDate
	transaction.Note = note
	transaction.UpdatedAt = time.Now()

	if err := repository.UpdateTransaction(transaction, db); err != nil {
		return nil, err
	}

	// Update new account balance
	newAccount, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}
	if transactionType == models.TransactionTypeIncome {
		newAccount.Balance += amount
	} else {
		newAccount.Balance -= amount
	}
	if err := repository.UpdateAccount(newAccount, db); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func DeleteTransaction(id uuid.UUID, db *sql.DB) error {
	transaction, err := repository.GetTransactionByID(id, db)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback on error, commit on success

	account, err := repository.GetAccountByID(transaction.AccountID, db)
	if err != nil {
		return err
	}

	if transaction.Type == models.TransactionTypeIncome {
		account.Balance -= transaction.Amount
	} else {
		account.Balance += transaction.Amount
	}

	if err := repository.UpdateAccount(account, db); err != nil {
		return err
	}

	if err := repository.DeleteTransaction(id, db); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetAggregateData(userID uuid.UUID, startDate string, endDate string, db *sql.DB) (map[string]interface{}, error) {
	return repository.GetAggregateDataByUserID(userID, startDate, endDate, db)
}



func GetSpendingByCategory(userID uuid.UUID, db *sql.DB) ([]map[string]interface{}, error) {
	return repository.GetSpendingByCategory(userID, db)
}
