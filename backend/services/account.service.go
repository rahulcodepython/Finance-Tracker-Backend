package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func CreateAccount(userID uuid.UUID, name string, accountType models.AccountType, balance float64, db *sql.DB) (*models.Account, error) {
	account := &models.Account{
		ID:       uuid.New(),
		UserID:   userID,
		Name:     name,
		Type:     accountType,
		Balance:  balance,
		IsActive: true,
	}

	err := repository.CreateAccount(account, db)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func GetAccounts(userID uuid.UUID, db *sql.DB) ([]models.Account, error) {
	return repository.GetAccountsByUserID(userID, db)
}

func CheckAccountExistsById(id uuid.UUID, db *sql.DB) (bool, error) {
	account, err := repository.GetAccountByID(id, db)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if account != nil {
		return true, nil
	}

	return false, nil
}

func UpdateAccount(id uuid.UUID, name string, accountType models.AccountType, isActive bool, db *sql.DB) (*models.Account, error) {
	account, err := repository.GetAccountByID(id, db)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, sql.ErrNoRows
	}

	account.Name = name
	account.Type = accountType
	account.IsActive = isActive
	account.UpdatedAt = time.Now().In(utils.LOC)

	err = repository.UpdateAccount(account, db)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func DeleteAccount(id uuid.UUID, db *sql.DB) error {
	account, err := repository.GetAccountByID(id, db)
	if err != nil {
		return err
	}

	if account == nil {
		return sql.ErrNoRows
	}

	return repository.DeleteAccount(id, db)
}

func GetTotalBalance(userID uuid.UUID, db *sql.DB) (float64, error) {
	accounts, err := repository.GetAccountsByUserID(userID, db)
	if err != nil {
		return 0, err
	}

	var totalBalance float64
	for _, account := range accounts {
		if account.IsActive {
			totalBalance += account.Balance
		}
	}

	return totalBalance, nil
}
