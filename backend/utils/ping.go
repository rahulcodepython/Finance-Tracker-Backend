// Package utils provides utility functions for the application.
package utils

import (
	"database/sql"
	"log"
)

// Ping checks the connection to the database and logs the status.
func Ping(db *sql.DB) error {
	// Ping the database to verify the connection is alive.
	if err := db.Ping(); err != nil {
		// If the ping fails, log the error and return it.
		log.Println("Unable to ping database:", err)
		return err
	}

	// If the ping is successful, log a success message.
	log.Println("Database is healthy.")
	return nil
}