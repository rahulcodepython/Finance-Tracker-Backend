package utils

import (
	"database/sql"
	"log"
)

func Ping(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		// If the ping fails, log the error and return it.
		log.Println("Unable to ping database:", err)
		return err // <-- Return the error
	}

	// If the ping is successful, a success message is logged.
	log.Println("Database is healthy.")
	return nil // <-- Return nil for success
}
