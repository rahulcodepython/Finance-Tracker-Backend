package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/database"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func DeserializeUser(c *fiber.Ctx) error {
	cfg := c.Locals("cfg").(*config.Config)

	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	}

	if tokenString == "" {
		return utils.UnauthorizedAccess(c, nil, "You are not logged in")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}
		return []byte(cfg.JWT.JWTSecret), nil
	})

	if err != nil {
		return utils.UnauthorizedAccess(c, err, "Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)

		db := database.DB

		jwtToken, err := repository.GetJwtTokenByToken(db, tokenString)
		if err != nil {
			return utils.UnauthorizedAccess(c, err, "Invalid token")
		}

		if jwtToken == nil {
			return utils.UnauthorizedAccess(c, err, "Invalid token")
		} else if jwtToken.ExpiresAt.Before(time.Now()) {
			err := repository.DeleteJwtToken(db, tokenString)
			if err != nil {
				return utils.UnauthorizedAccess(c, err, "Invalid token")
			}
		}

		c.Locals("user_id", userID)
		return c.Next()
	}

	return utils.UnauthorizedAccess(c, err, "Invalid token")

}
