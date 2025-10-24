package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func Migrate(db *sql.DB) {
	sqlFilePath := "migrations/schema.sql"
	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Printf("Error reading SQL file %s: %v", sqlFilePath, err)
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
		err := CommitExec(trimmedQuery, db)
		if err != nil {
			// If a statement fails, we log the error and the problematic query, then stop.
			log.Printf("Failed to execute command Error: %v for sql query %s", err, trimmedQuery)
		}
	}

	log.Println("Database schema setup complete.")
}
