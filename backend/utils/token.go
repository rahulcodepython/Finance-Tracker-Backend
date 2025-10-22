package utils

import (
	"time"

	"github.com/rahulcodepython/finance-tracker-backend/backend/config"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID string, cfg *config.Config) (string, time.Time, error) {
	secret := cfg.JWT.JWTSecret
	expiresIn := cfg.JWT.JWTExpiresIn

	expirationTime, err := time.ParseDuration(expiresIn)
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(expirationTime)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
