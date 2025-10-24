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

func CheckUserExistsByEmail(email string, db *sql.DB) (bool, error) {
	existingUser, err := repository.GetUserByEmail(email, db)
	if err != nil {
		return false, err
	}

	if existingUser != nil {
		return true, nil
	}

	return false, nil
}

func Register(name, email, password string, db *sql.DB, cfg *config.Config) (*models.User, string, error) {
	exists, err := CheckUserExistsByEmail(email, db)

	if err != nil {
		return nil, "", err
	}

	if exists {
		return nil, "", errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		Provider:  models.AuthProviderEmail,
		CreatedAt: time.Now().In(utils.LOC),
	}

	err = repository.CreateUser(user, db)
	if err != nil {
		return nil, "", err
	}

	// Log the registration
	go CreateLog(user.ID, "User registered", db)

	token, expiresAt, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	jwtToken := models.JwtToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().In(utils.LOC),
	}

	repository.CreateJwtToken(db, &jwtToken)

	// Log the registration
	go CreateLog(user.ID, "JWT token created", db)

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

	// Log the login
	go CreateLog(user.ID, "User logged in", db)

	jwtToken, err := repository.GetJwtTokenByUserID(db, user.ID)
	if err != nil {
		return nil, "", err
	}

	if jwtToken != nil && jwtToken.ExpiresAt.After(time.Now().In(utils.LOC)) {
		return user, jwtToken.Token, nil
	}

	if jwtToken != nil && jwtToken.ExpiresAt.Before(time.Now().In(utils.LOC)) {
		err := repository.DeleteJwtTokenByUserID(db, user.ID)
		if err != nil {
			return nil, "", err
		}

		// Log the deletion
		go CreateLog(user.ID, "JWT token deleted", db)
	}

	token, expiresAt, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	newJwtToken := models.JwtToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().In(utils.LOC),
	}

	repository.CreateJwtToken(db, &newJwtToken)

	// Log the registration
	go CreateLog(user.ID, "JWT token created", db)

	return user, token, nil
}

func ChangePassword(userID uuid.UUID, currentPassword, newPassword string, db *sql.DB) error {
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

	err = repository.UpdateUser(user, db)
	if err != nil {
		return err
	}

	// Log the password change
	go CreateLog(user.ID, "User changed password", db)

	return nil
}

func GetProfile(userID uuid.UUID, db *sql.DB) (*models.User, error) {
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
			CreatedAt: time.Now().In(utils.LOC),
		}

		if err := repository.CreateUser(user, db); err != nil {
			return nil, "", err
		}
		// Log the registration
		go CreateLog(user.ID, "User registered with Google", db)
	}

	// Log the login
	go CreateLog(user.ID, "User logged in with Google", db)

	jwtToken, err := repository.GetJwtTokenByUserID(db, user.ID)
	if err != nil {
		return nil, "", err
	}

	if jwtToken != nil && jwtToken.ExpiresAt.After(time.Now().In(utils.LOC)) {
		return user, jwtToken.Token, nil
	}

	if jwtToken != nil && jwtToken.ExpiresAt.Before(time.Now().In(utils.LOC)) {
		err := repository.DeleteJwtTokenByUserID(db, user.ID)
		if err != nil {
			return nil, "", err
		}

		// Log the deletion
		go CreateLog(user.ID, "JWT token deleted", db)
	}

	token, expiresAt, err := utils.GenerateToken(user.ID.String(), cfg)
	if err != nil {
		return nil, "", err
	}

	newJwtToken := models.JwtToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().In(utils.LOC),
	}

	repository.CreateJwtToken(db, &newJwtToken)

	// Log the registration
	go CreateLog(user.ID, "JWT token created", db)

	return user, token, nil
}
