// Package middleware provides middleware functions for the Fiber web framework.
package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// DeserializeUser is a middleware function that deserializes the user from a JWT token.
// It validates the token, checks for its existence and expiration, and sets the user ID in the context.
func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string

	// Get the configuration from the context.
	cfg := c.Locals("cfg").(*config.Config)
	// Get the Authorization header.
	authorization := c.Get("Authorization")

	// Check if the Authorization header is in the correct format.
	if !strings.HasPrefix(authorization, "Bearer ") {
		return utils.UnauthorizedAccess(c, errors.New("invalid token format"), "Unauthorized Access")
	}

	// Extract the token from the Authorization header.
	tokenString = strings.TrimPrefix(authorization, "Bearer ")

	// Check if the token is empty.
	if tokenString == "" {
		return utils.UnauthorizedAccess(c, errors.New("empty token"), "Unauthorized Access")
	}

	// Parse the JWT token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}
		// Return the JWT secret.
		return []byte(cfg.JWT.JWTSecret), nil
	})

	// Check for parsing errors.
	if err != nil {
		return utils.UnauthorizedAccess(c, err, "Invalid token")
	}

	// Check if the token is valid and contains claims.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user ID from the claims.
		userID := claims["user_id"].(string)

		// Get the database connection from the context.
		db := database.DB

		// Check if the token exists in the database.
		jwtToken, err := repository.GetJwtTokenByToken(db, tokenString)
		if err != nil {
			return utils.UnauthorizedAccess(c, err, "Invalid token")
		}

		// If the token does not exist or has expired, return an unauthorized error.
		if jwtToken == nil {
			return utils.UnauthorizedAccess(c, err, "Invalid token")
		} else if jwtToken.ExpiresAt.Before(time.Now().In(utils.LOC)) {
			// If the token has expired, delete it from the database.
			err := repository.DeleteJwtToken(db, tokenString)
			if err != nil {
				return utils.UnauthorizedAccess(c, err, "Invalid token")
			}
		}

		// Set the user ID in the context for later use.
		c.Locals("user_id", userID)
		// Continue to the next middleware or handler.
		return c.Next()
	}

	// Return an unauthorized error if the token is invalid.
	return utils.UnauthorizedAccess(c, err, "Invalid token")

}
