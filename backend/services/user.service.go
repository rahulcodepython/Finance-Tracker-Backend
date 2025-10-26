// Package services provides business logic for the application.
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

// CheckUserExistsByEmail checks if a user exists by their email address.
func CheckUserExistsByEmail(email string, db *sql.DB) (bool, error) {
	// Get the user from the database.
	existingUser, err := repository.GetUserByEmail(email, db)
	if err != nil {
		return false, err
	}

	// If the user is not nil, they exist.
	if existingUser != nil {
		return true, nil
	}

	return false, nil
}

// Register creates a new user and returns the user and a JWT token.
func Register(name, email, password string, db *sql.DB, cfg *config.Config) (*models.User, string, error) {
	// Check if a user with the given email already exists.
	exists, err := CheckUserExistsByEmail(email, db)
	if err != nil {
		return nil, "", err
	}

	if exists {
		return nil, "", errors.New("user already exists")
	}

	// Hash the user's password.
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// Create a new User model.
	user := &models.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		Provider:  models.AuthProviderEmail,
		CreatedAt: time.Now().In(utils.LOC),
	}

	// Create the user in the database.
	err = repository.CreateUser(user, db)
	if err != nil {
		return nil, "", err
	}

	// Create a log entry for the user registration.
	go CreateLog(user.ID, "User registered", db)

	// Generate a JWT token for the new user.
	token, expiresAt, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	// Create a new JwtToken model.
	jwtToken := models.JwtToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().In(utils.LOC),
	}

	// Create the JWT token in the database.
	repository.CreateJwtToken(db, &jwtToken)

	// Create a log entry for the JWT token creation.
	go CreateLog(user.ID, "JWT token created", db)

	return user, token, nil
}

// Login authenticates a user and returns the user and a JWT token.
func Login(email, password string, db *sql.DB, cfg *config.Config) (*models.User, string, error) {
	// Get the user from the database by email.
	user, err := repository.GetUserByEmail(email, db)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Check if the provided password matches the stored hash.
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	// Create a log entry for the user login.
	go CreateLog(user.ID, "User logged in", db)

	// Get the user's JWT token from the database.
	jwtToken, err := repository.GetJwtTokenByUserID(db, user.ID)
	if err != nil {
		return nil, "", err
	}

	// If a valid token exists, return it.
	if jwtToken != nil && jwtToken.ExpiresAt.After(time.Now().In(utils.LOC)) {
		return user, jwtToken.Token, nil
	}

	// If an expired token exists, delete it.
	if jwtToken != nil && jwtToken.ExpiresAt.Before(time.Now().In(utils.LOC)) {
		err := repository.DeleteJwtTokenByUserID(db, user.ID)
		if err != nil {
			return nil, "", err
		}

		// Create a log entry for the token deletion.
		go CreateLog(user.ID, "JWT token deleted", db)
	}

	// Generate a new JWT token.
	token, expiresAt, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	// Create a new JwtToken model.
	newJwtToken := models.JwtToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().In(utils.LOC),
	}

	// Create the new JWT token in the database.
	repository.CreateJwtToken(db, &newJwtToken)

	// Create a log entry for the new token creation.
	go CreateLog(user.ID, "JWT token created", db)

	return user, token, nil
}

// ChangePassword changes a user's password.
func ChangePassword(userID uuid.UUID, currentPassword, newPassword string, db *sql.DB) error {
	// Get the user from the database.
	user, err := repository.GetUserByID(userID, db)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if the current password is correct.
	if !utils.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("invalid current password")
	}

	// Hash the new password.
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update the user's password.
	user.Password = hashedPassword

	// Update the user in the database.
	err = repository.UpdateUser(user, db)
	if err != nil {
		return err
	}

	// Create a log entry for the password change.
	go CreateLog(user.ID, "User changed password", db)

	return nil
}

// GetProfile retrieves a user's profile information.
func GetProfile(userID uuid.UUID, db *sql.DB) (*models.User, error) {
	return repository.GetUserByID(userID, db)
}

// GoogleLogin handles user login or registration via Google OAuth.
func GoogleLogin(email, fullName string, db *sql.DB, cfg *config.Config) (*models.User, string, error) {
	// Get the user from the database by email.
	user, err := repository.GetUserByEmail(email, db)
	if err != nil {
		// If the user does not exist, create a new user.
		user = &models.User{
			ID:        uuid.New(),
			Name:      fullName,
			Email:     email,
			Provider:  models.AuthProviderGoogle,
			CreatedAt: time.Now().In(utils.LOC),
		}

		if err := repository.CreateUser(user, db); err != nil {
			return nil, "", err
		}
		// Create a log entry for the user registration.
		go CreateLog(user.ID, "User registered with Google", db)
	}

	// Create a log entry for the user login.
	go CreateLog(user.ID, "User logged in with Google", db)

	// Get the user's JWT token from the database.
	jwtToken, err := repository.GetJwtTokenByUserID(db, user.ID)
	if err != nil {
		return nil, "", err
	}

	// If a valid token exists, return it.
	if jwtToken != nil && jwtToken.ExpiresAt.After(time.Now().In(utils.LOC)) {
		return user, jwtToken.Token, nil
	}

	// If an expired token exists, delete it.
	if jwtToken != nil && jwtToken.ExpiresAt.Before(time.Now().In(utils.LOC)) {
		err := repository.DeleteJwtTokenByUserID(db, user.ID)
		if err != nil {
			return nil, "", err
		}

		// Create a log entry for the token deletion.
		go CreateLog(user.ID, "JWT token deleted", db)
	}

	// Generate a new JWT token.
	token, expiresAt, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	// Create a new JwtToken model.
	newJwtToken := models.JwtToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().In(utils.LOC),
	}

	// Create the new JWT token in the database.
	repository.CreateJwtToken(db, &newJwtToken)

	// Create a log entry for the new token creation.
	go CreateLog(user.ID, "JWT token created", db)

	return user, token, nil
}