package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	cfg := c.Locals("cfg").(*config.Config)

	if len(authorization) > 7 && authorization[:7] == "Bearer " {
		tokenString = authorization[7:]
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
		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}

	return utils.UnauthorizedAccess(c, err, "Invalid token")
}
