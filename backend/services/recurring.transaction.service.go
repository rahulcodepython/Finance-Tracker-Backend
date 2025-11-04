// Package services provides business logic for the application.
package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/serializers"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateRecurringTransaction creates a new recurring transaction for a user.
func CreateRecurringTransaction(userID uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, note sql.NullString, recurringFrequency models.RecurringFrequency, recurringDate int, db *sql.DB) (*serializers.RecurringTransactionResponse, error) {
	// Get the category to determine the transaction type.
	category, err := repository.GetCategoryByID(categoryID, db)
	if err != nil {
		return nil, err
	}

	// If the category does not exist, return an error.
	if category == nil {
		return nil, sql.ErrNoRows
	}

	// Set the transaction type based on the category type.
	transactionType := models.TransactionType(category.Type)

	// Get the account to ensure it exists.
	account, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}

	// If the account does not exist, return an error.
	if account == nil {
		return nil, sql.ErrNoRows
	}

	// If a budget is provided, ensure it exists.
	if budgetID.Valid {
		budget, err := repository.GetBudgetByID(budgetID.UUID, db)
		if err != nil {
			return nil, err
		}

		if budget == nil {
			return nil, sql.ErrNoRows
		}
	}

	// Create a new RecurringTransaction model.
	recurringTransaction := &models.RecurringTransaction{
		ID:                 uuid.New(),
		UserID:             userID,
		AccountID:          accountID,
		CategoryID:         categoryID,
		BudgetID:           budgetID,
		Description:        description,
		Amount:             amount,
		Type:               transactionType,
		Note:               note,
		RecurringFrequency: recurringFrequency,
		RecurringDate:      recurringDate,
		CreatedAt:          time.Now().In(utils.LOC),
		UpdatedAt:          time.Now().In(utils.LOC),
	}

	// Create the recurring transaction in the database.
	if err := repository.CreateRecurringTransaction(recurringTransaction, db); err != nil {
		return nil, err
	}

	// Create a log entry for the recurring transaction creation.
	go CreateLog(userID, fmt.Sprintf("New recurring transaction '%s' created", recurringTransaction.Description), db)

	var recurringTransactionResponse serializers.RecurringTransactionResponse
	recurringTransactionResponse.ID = recurringTransaction.ID
	recurringTransactionResponse.UserID = recurringTransaction.UserID
	recurringTransactionResponse.Description = recurringTransaction.Description
	recurringTransactionResponse.Amount = recurringTransaction.Amount
	recurringTransactionResponse.Type = recurringTransaction.Type
	recurringTransactionResponse.Note = recurringTransaction.Note
	recurringTransactionResponse.RecurringFrequency = recurringTransaction.RecurringFrequency
	recurringTransactionResponse.RecurringDate = recurringTransaction.RecurringDate
	recurringTransactionResponse.CreatedAt = recurringTransaction.CreatedAt
	recurringTransactionResponse.UpdatedAt = recurringTransaction.UpdatedAt
	recurringTransactionResponse.Account.UUID = recurringTransaction.AccountID
	recurringTransactionResponse.Category.UUID = recurringTransaction.CategoryID
	if recurringTransaction.BudgetID.Valid {
		recurringTransactionResponse.Budget.UUID = recurringTransaction.BudgetID.UUID
	}

	return &recurringTransactionResponse, nil
}

// GetRecurringTransactions retrieves all recurring transactions for a user.
func GetRecurringTransactions(userID uuid.UUID, db *sql.DB) ([]serializers.RecurringTransactionResponse, error) {
	return repository.GetRecurringTransactionsByUserID(userID, db)
}

// UpdateRecurringTransaction updates an existing recurring transaction.
func UpdateRecurringTransaction(id uuid.UUID, accountID uuid.UUID, categoryID uuid.UUID, budgetID uuid.NullUUID, description string, amount float64, note sql.NullString, recurringFrequency models.RecurringFrequency, recurringDate int, db *sql.DB) (*serializers.RecurringTransactionResponse, error) {
	// Get the recurring transaction from the database.
	recurringTransaction, err := repository.GetRecurringTransactionByID(id, db)
	if err != nil {
		return nil, err
	}

	// If the recurring transaction does not exist, return an error.
	if recurringTransaction == nil {
		return nil, sql.ErrNoRows
	}

	// Get the category to determine the transaction type.
	category, err := repository.GetCategoryByID(categoryID, db)
	if err != nil {
		return nil, err
	}

	// If the category does not exist, return an error.
	if category == nil {
		return nil, sql.ErrNoRows
	}

	// Set the transaction type based on the category type.
	transactionType := models.TransactionType(category.Type)

	// Get the account to ensure it exists.
	account, err := repository.GetAccountByID(accountID, db)
	if err != nil {
		return nil, err
	}

	// If the account does not exist, return an error.
	if account == nil {
		return nil, sql.ErrNoRows
	}

	// If a budget is provided, ensure it exists.
	if budgetID.Valid {
		budget, err := repository.GetBudgetByID(budgetID.UUID, db)
		if err != nil {
			return nil, err
		}

		if budget == nil {
			return nil, sql.ErrNoRows
		}
	}

	// Update the recurring transaction fields.
	recurringTransaction.AccountID = accountID
	recurringTransaction.CategoryID = categoryID
	recurringTransaction.BudgetID = budgetID
	recurringTransaction.Description = description
	recurringTransaction.Amount = amount
	recurringTransaction.Type = transactionType
	recurringTransaction.Note = note
	recurringTransaction.RecurringFrequency = recurringFrequency
	recurringTransaction.RecurringDate = recurringDate
	recurringTransaction.UpdatedAt = time.Now().In(utils.LOC)

	// Update the recurring transaction in the database.
	if err := repository.UpdateRecurringTransaction(recurringTransaction, db); err != nil {
		return nil, err
	}

	// Create a log entry for the recurring transaction update.
	go CreateLog(recurringTransaction.UserID, fmt.Sprintf("Recurring transaction '%s' updated", recurringTransaction.Description), db)

	var recurringTransactionResponse serializers.RecurringTransactionResponse
	recurringTransactionResponse.ID = recurringTransaction.ID
	recurringTransactionResponse.UserID = recurringTransaction.UserID
	recurringTransactionResponse.Description = recurringTransaction.Description
	recurringTransactionResponse.Amount = recurringTransaction.Amount
	recurringTransactionResponse.Type = recurringTransaction.Type
	recurringTransactionResponse.Note = recurringTransaction.Note
	recurringTransactionResponse.RecurringFrequency = recurringTransaction.RecurringFrequency
	recurringTransactionResponse.RecurringDate = recurringTransaction.RecurringDate
	recurringTransactionResponse.CreatedAt = recurringTransaction.CreatedAt
	recurringTransactionResponse.UpdatedAt = recurringTransaction.UpdatedAt
	recurringTransactionResponse.Account.UUID = recurringTransaction.AccountID
	recurringTransactionResponse.Category.UUID = recurringTransaction.CategoryID
	if recurringTransaction.BudgetID.Valid {
		recurringTransactionResponse.Budget.UUID = recurringTransaction.BudgetID.UUID
	}

	return &recurringTransactionResponse, nil
}

// DeleteRecurringTransaction deletes a recurring transaction.
func DeleteRecurringTransaction(id uuid.UUID, db *sql.DB) error {
	// Get the recurring transaction from the database.
	recurringTransaction, err := repository.GetRecurringTransactionByID(id, db)
	if err != nil {
		return err
	}

	// If the recurring transaction does not exist, return an error.
	if recurringTransaction == nil {
		return sql.ErrNoRows
	}

	// Delete the recurring transaction from the database.
	err = repository.DeleteRecurringTransaction(id, db)
	if err != nil {
		return err
	}

	// Create a log entry for the recurring transaction deletion.
	go CreateLog(recurringTransaction.UserID, fmt.Sprintf("Recurring transaction '%s' removed", recurringTransaction.Description), db)

	return nil
}
