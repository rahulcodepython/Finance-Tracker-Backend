// Package utils provides utility functions for the application.
package utils

import (
	"github.com/gofiber/fiber/v2" // Web framework for Go
)

// Response defines the structure of a standard API response.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// InternalServerError sends a 500 Internal Server Error response.
func InternalServerError(c *fiber.Ctx, err error, message string) error {
	// Use a default message if none is provided.
	if message == "" {
		message = "Internal Server Error"
	}

	// Get the error message string.
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}

	// Send the JSON response.
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Message: message,
		Error:   errMessage,
	})
}

// UnauthorizedAccess sends a 401 Unauthorized response.
func UnauthorizedAccess(c *fiber.Ctx, err error, message string) error {
	// Use a default message if none is provided.
	if message == "" {
		message = "Unauthorized Access"
	}

	// Get the error message string.
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}

	// Send the JSON response.
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Message: message,
		Error:   errMessage,
	})
}

// NotFound sends a 404 Not Found response.
func NotFound(c *fiber.Ctx, err error, message string) error {
	// Use a default message if none is provided.
	if message == "" {
		message = "Not Found"
	}

	// Get the error message string.
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}

	// Send the JSON response.
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Message: message,
		Error:   errMessage,
	})
}

// BadResponse sends a 400 Bad Request response.
func BadResponse(c *fiber.Ctx, err error, message string) error {
	// Use a default message if none is provided.
	if message == "" {
		message = "Bad Request"
	}

	// Get the error message string.
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}

	// Send the JSON response.
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Message: message,
		Error:   errMessage,
	})
}

// OKResponse sends a 200 OK response.
func OKResponse(c *fiber.Ctx, message string, data interface{}) error {
	// Send the JSON response.
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// OKCreatedResponse sends a 201 Created response.
func OKCreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	// Send the JSON response.
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// TooManyRequests sends a 429 Too Many Requests response.
func TooManyRequests(c *fiber.Ctx, message string) error {
	// Send the JSON response.
	return c.Status(fiber.StatusTooManyRequests).JSON(Response{
		Success: false,
		Message: message,
	})
}