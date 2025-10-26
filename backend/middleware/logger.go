// Package middleware provides middleware functions for the Fiber web framework.
package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2" // Web framework for Go
	"github.com/google/uuid"
)

// Logger is a middleware that logs HTTP requests in a detailed format.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start a timer to measure the request processing time.
		start := time.Now()

		// Process the request by calling the next middleware or handler.
		err := c.Next()

		// Stop the timer.
		stop := time.Now()

		// Determine the user for logging purposes.
		user := "anonymous"
		userIDVal := c.Locals("user_id")
		if userIDVal != nil {
			userID, _ := uuid.Parse(userIDVal.(string))
			user = userID.String()
		}

		// Get the request ID from the header or generate a new one.
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Create the log message with various request and response details.
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

		// Append any errors to the log message.
		if err != nil {
			log = fmt.Sprintf("%s\nError: %v", log, err)
		}

		// Print the log message to the console.
		fmt.Println(log)

		return err
	}
}