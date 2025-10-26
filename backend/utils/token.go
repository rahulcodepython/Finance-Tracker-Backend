// Package utils provides utility functions for the application.
package utils

import (
	"time"

	"github.com/rahulcodepython/finance-tracker-backend/backend/config"

	"github.com/golang-jwt/jwt/v4" // Used for generating JWT tokens
)

// GenerateToken generates a new JWT token for a given user ID.
func GenerateToken(userID string, cfg *config.Config) (string, time.Time, error) {
	// Get the JWT secret and expiration duration from the configuration.
	secret := cfg.JWT.JWTSecret
	expiresIn := cfg.JWT.JWTExpiresIn

	// Parse the expiration duration string.
	expirationTime, err := time.ParseDuration(expiresIn)
	if err != nil {
		return "", time.Time{}, err
	}

	// Calculate the token expiration time.
	expiresAt := time.Now().In(LOC).Add(expirationTime)
	// Create the JWT claims, including the user ID and expiration time.
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	}

	// Create a new JWT token with the specified signing method and claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key.
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}