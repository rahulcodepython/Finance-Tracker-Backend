// Package utils provides utility functions for the application.
package utils

import (
	"log"
	"time"
)

// LOC is a global variable that holds the loaded timezone.
var LOC *time.Location

// LoadTimezone loads the Asia/Kolkata timezone and stores it in the LOC variable.
func LoadTimezone() {
	var err error
	// Load the specified timezone.
	LOC, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// If loading fails, log the error.
		log.Printf("Failed to load timezone: %v", err)
	}
}