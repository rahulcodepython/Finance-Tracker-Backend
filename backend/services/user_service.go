package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func Register(fullName, email, password string, db *sql.DB, cfg *config.Config) (*models.User, string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		ID:        uuid.New(),
		Name:      fullName,
		Email:     email,
		Password:  hashedPassword,
		Provider:  models.AuthProviderEmail,
		CreatedAt: time.Now(),
	}

	err = repository.CreateUser(user, db)
	if err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func Login(email, password string, db *sql.DB, cfg *config.Config) (*models.User, string, error) {
	user, err := repository.GetUserByEmail(email, db)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func ChangePassword(userID, currentPassword, newPassword string, db *sql.DB) error {
	user, err := repository.GetUserByID(userID, db)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("invalid current password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return repository.UpdateUser(user, db)
}

func GetProfile(userID string, db *sql.DB) (*models.User, error) {
	return repository.GetUserByID(userID, db)
}

func GoogleLogin(email, fullName string, db *sql.DB, cfg *config.Config) (*models.User, string, error) {
	user, err := repository.GetUserByEmail(email, db)
	if err != nil {
		// User does not exist, create a new user
		user = &models.User{
			ID:        uuid.New(),
			Name:      fullName,
			Email:     email,
			Provider:  models.AuthProviderGoogle,
			CreatedAt: time.Now(),
		}

		if err := repository.CreateUser(user, db); err != nil {
			return nil, "", err
		}
	}

	token, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
