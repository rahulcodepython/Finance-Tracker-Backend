package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

var DB *sql.DB

func Connect(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kolkata",
		cfg.Database.DBHost,
		cfg.Database.DBUser,
		cfg.Database.DBPassword,
		cfg.Database.DBName,
		cfg.Database.DBPort,
		cfg.Database.DBSSMode,
	)

	var err error

	DB, err = sql.Open("postgres", dsn)
	// This checks if an error occurred while opening the database connection.
	if err != nil {
		// If an error occurs, a message is logged.
		log.Println("Unable to connect with database")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}

	// PingDB() is called to check if the database connection is alive.
	utils.Ping(DB)

	// The database connection is returned.
	return DB
}
