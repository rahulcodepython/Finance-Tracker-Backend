// Package database provides functionality for connecting to and interacting with the database.
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// DB is a global variable that holds the database connection.
var DB *sql.DB

// Connect establishes a connection to the database using the provided configuration.
func Connect(cfg *config.Config) *sql.DB {
	// Create the data source name (DSN) string from the configuration.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kolkata",
		cfg.Database.DBHost,
		cfg.Database.DBUser,
		cfg.Database.DBPassword,
		cfg.Database.DBName,
		cfg.Database.DBPort,
		cfg.Database.DBSSMode,
	)

	var err error

	// Open a new database connection.
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Println("Unable to connect with database")
		log.Fatal(err)
	}

	// Ping the database to verify the connection.
	utils.Ping(DB)

	// Return the database connection.
	return DB
}