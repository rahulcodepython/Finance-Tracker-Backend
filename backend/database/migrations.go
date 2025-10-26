// Package database provides functionality for connecting to and interacting with the database.
package database

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color" // Used for colorizing console output
	_ "github.com/lib/pq"    // PostgreSQL driver
)

// Migrate reads and executes SQL commands from the schema.sql file to set up the database.
func Migrate(db *sql.DB) {
	// Define the path to the SQL migration file.
	sqlFilePath := "migrations/schema.sql"

	// Read the content of the SQL file.
	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("❌ Failed to read SQL file %s: %v", sqlFilePath, err)
	}

	log.Println("🚀 Starting database migration...")
	start := time.Now()

	// Split the file content into individual SQL queries.
	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		// Handle ENUM type creation safely by checking for existence first.
		if isEnumCreate(query) {
			typeName := extractEnumName(query)
			if typeExists(db, typeName) {
				log.Printf("⚡ Enum type %s already exists. Skipping.", typeName)
				continue
			}
		}

		// Execute the SQL query.
		if _, err := db.Exec(query); err != nil {
			log.Printf("❌ ERROR: Failed executing SQL → %v\n    ↳ %s", err, previewQuery(query))
			continue
		}

		// Log a success message for the executed query.
		logSuccess(query)
	}

	elapsed := time.Since(start)
	green := color.New(color.FgGreen).SprintFunc()
	log.Printf("%s Database migration completed successfully in %s\n", green("✅ DONE:"), elapsed.Round(time.Millisecond))
}

// isEnumCreate checks if a query is a 'CREATE TYPE ... AS ENUM' statement.
func isEnumCreate(query string) bool {
	matched, _ := regexp.MatchString(`(?i)^CREATE\s+TYPE\s+\w+\s+AS\s+ENUM`, query)
	return matched
}

// extractEnumName extracts the name of the ENUM type from a 'CREATE TYPE' query.
func extractEnumName(query string) string {
	re := regexp.MustCompile(`(?i)^CREATE\s+TYPE\s+(\w+)\s+AS\s+ENUM`)
	matches := re.FindStringSubmatch(query)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown_enum"
}

// typeExists checks if a given type already exists in the database.
func typeExists(db *sql.DB, typeName string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname=$1);`
	if err := db.QueryRow(query, typeName).Scan(&exists); err != nil {
		log.Printf("❌ ERROR checking enum existence %s: %v", typeName, err)
		return false
	}
	return exists
}

// logSuccess logs a formatted success message based on the type of SQL query executed.
func logSuccess(query string) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	upper := strings.ToUpper(strings.TrimSpace(query))

	switch {
	case strings.HasPrefix(upper, "CREATE TABLE"):
		table := extractName(query, "CREATE TABLE IF NOT EXISTS", "CREATE TABLE")
		log.Printf("%s Table %s created successfully.", green("✅"), yellow(table))
	case strings.HasPrefix(upper, "DROP TABLE"):
		table := extractName(query, "DROP TABLE IF EXISTS", "DROP TABLE")
		log.Printf("%s Table %s dropped successfully.", green("🗑️"), yellow(table))
	case strings.HasPrefix(upper, "CREATE TYPE"):
		typ := extractEnumName(query)
		log.Printf("%s Type %s created successfully.", green("✅"), yellow(typ))
	case strings.HasPrefix(upper, "DROP TYPE"):
		typ := extractName(query, "DROP TYPE IF EXISTS", "DROP TYPE")
		log.Printf("%s Type %s dropped successfully.", green("🗑️"), yellow(typ))
	case strings.HasPrefix(upper, "CREATE INDEX"):
		idx := extractName(query, "CREATE INDEX IF NOT EXISTS", "CREATE INDEX")
		log.Printf("%s Index %s created successfully.", green("📈"), yellow(idx))
	case strings.HasPrefix(upper, "DROP INDEX"):
		idx := extractName(query, "DROP INDEX IF EXISTS", "DROP INDEX")
		log.Printf("%s Index %s dropped successfully.", green("🗑️"), yellow(idx))
	default:
		log.Printf("%s Executed SQL successfully.", green("✔️"))
	}
}

// extractName extracts the name of the table, index, or type from a SQL query.
func extractName(query string, prefixes ...string) string {
	upper := strings.ToUpper(query)
	for _, prefix := range prefixes {
		p := strings.ToUpper(prefix)
		if strings.HasPrefix(upper, p) {
			rest := strings.TrimSpace(query[len(p):])
			parts := strings.Fields(rest)
			if len(parts) > 0 {
				return strings.Trim(parts[0], `"`)
			}
		}
	}
	return "unknown"
}

// previewQuery returns a short preview of a SQL query for logging purposes.
func previewQuery(query string) string {
	trimmed := strings.ReplaceAll(strings.TrimSpace(query), "\n", " ")
	if len(trimmed) > 120 {
		return trimmed[:120] + "..."
	}
	return trimmed
}
