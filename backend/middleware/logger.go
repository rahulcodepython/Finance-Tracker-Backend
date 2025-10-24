// This file defines a middleware for logging HTTP requests.
package middleware

// "github.com/gofiber/fiber/v2" is a web framework for Go. It is used here to create middleware.
import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Logger is a middleware that logs HTTP requests.
// It returns a Fiber handler.

// @return fiber.Handler - The Fiber handler.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Stop timer
		stop := time.Now()

		// Get user from context
		user := "anonymous"
		userIDVal := c.Locals("user_id")
		if userIDVal != nil {
			userID, _ := uuid.Parse(userIDVal.(string))
			user = userID.String()
		}

		// Get request ID
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Log format
		log := fmt.Sprintf("[%s] IP: %s - Method: %s - Path: %s - Status: %d - Duration: %s - User: %s - Protocol: %s - Request ID: %s - User Agent: %s",
			start.Format("2006-01-02 15:04:05"),
			c.IP(),
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			stop.Sub(start),
			user,
			c.Protocol(),
			requestID,
			c.Get("User-Agent"),
		)

		if err != nil {
			log = fmt.Sprintf("%s\nError: %v", log, err)
		}

		// Print log
		fmt.Println(log)

		return err
	}
}
