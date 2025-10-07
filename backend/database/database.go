package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rahulcodepython/finance-tracker-backend/backend/config"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func Connect(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Database.DBHost,
		cfg.Database.DBUser,
		cfg.Database.DBPassword,
		cfg.Database.DBName,
		cfg.Database.DBPort,
		cfg.Database.DBSSMode,
	)
	db, err := sql.Open("postgres", dsn)
	// This checks if an error occurred while opening the database connection.
	if err != nil {
		// If an error occurs, a message is logged.
		log.Println("Unable to connect with database")
		// The application is terminated with a fatal error.
		log.Fatal(err)
	}

	// PingDB() is called to check if the database connection is alive.
	utils.Ping(db)

	// The database connection is returned.
	return db
}

func CreateTables(db *sql.DB) {
	sqlFilePath := "migrations/schema.sql"
	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("Error reading SQL file %s: %v", sqlFilePath, err)
	}

	// 3. SPLIT INTO INDIVIDUAL STATEMENTS
	queries := strings.Split(string(content), ";")

	// 4. EXECUTE EACH STATEMENT USING db.Exec()
	for _, query := range queries {
		trimmedQuery := strings.TrimSpace(query)

		// Skip empty statements
		if trimmedQuery == "" {
			continue
		}

		// Execute the command. db.Exec() is for statements that don't return rows.
		_, err := db.Exec(trimmedQuery)
		if err != nil {
			// If a statement fails, we log the error and the problematic query, then stop.
			log.Fatalf("Failed to execute command:\n%s\n\nError: %v", trimmedQuery, err)
		}

		// Optional: Log successful execution of each statement
		// To avoid printing long CREATE TABLE statements, we can just print a snippet.
		getQuerySnippet := func(query string) string {
			words := strings.Fields(query)
			if len(words) > 5 {
				return strings.Join(words[:5], " ")
			}
			return query
		}
		log.Printf("Successfully executed statement starting with: \"%s...\"", getQuerySnippet(trimmedQuery))
	}

	fmt.Println("\nDatabase schema setup complete.")
}
