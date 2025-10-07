package utils

import (
	"time"

	"github.com/rahulcodepython/finance-tracker-backend/backend/config"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID string, cfg *config.Config) (string, error) {
	secret := cfg.JWT.JWTSecret
	expiresAt := cfg.JWT.JWTExpiresAt

	expirationTime, err := time.ParseDuration(expiresAt)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expirationTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
