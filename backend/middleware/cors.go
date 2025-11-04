// This file defines a middleware for handling Cross-Origin Resource Sharing (CORS).
package middleware

// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
import (
	"github.com/gofiber/fiber/v2"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"

	// "github.com/gofiber/fiber/v2/middleware/cors" is a middleware that provides CORS functionality.
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Cors is a middleware that handles CORS.
// It takes the application configuration as input and returns a Fiber handler.
//
// @param cfg *config.Config - The application configuration.
// @return fiber.Handler - The Fiber handler.
func Cors() fiber.Handler {
	// cors.New() returns a new CORS middleware with the specified configuration.
	return cors.New(cors.Config{
		// AllowOrigins is a list of origins that are allowed to make cross-origin requests.
		AllowOrigins: config.CFG.ServerConfig.ClientOrigin,
		// AllowHeaders is a list of headers that are allowed in cross-origin requests.
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		// AllowMethods is a list of methods are allowed in cross-origin requests.
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		// Next is a function that determines whether to skip this middleware.
		Next: func(c *fiber.Ctx) bool {
			// The middleware is skipped if the request is coming from the server itself.
			return c.IP() == config.CFG.ServerConfig.Host
		},
	})
}
