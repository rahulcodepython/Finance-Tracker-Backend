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

func ValidateToken(tokenString string, cfg *config.Config) (*jwt.Token, error) {
	secret := cfg.JWT.JWTSecret

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(secret), nil
	})
}
