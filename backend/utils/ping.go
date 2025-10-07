package utils

import (
	"database/sql"
	"log"
)

func Ping(db *sql.DB) {
	if err := db.Ping(); err != nil {
		// If the ping fails, a message is logged.
		log.Println("Unable to ping database")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}

	// If the ping is successful, a success message is logged.
	log.Println("Database is healthy.")
}
