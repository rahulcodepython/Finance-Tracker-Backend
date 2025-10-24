package database

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
)

// Migrate reads and executes SQL commands from a file dynamically
func Migrate(db *sql.DB) {
	sqlFilePath := "migrations/schema.sql"

	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("❌ Failed to read SQL file %s: %v", sqlFilePath, err)
	}

	log.Println("🚀 Starting database migration...")
	start := time.Now()

	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		// Handle dynamic ENUM creation safely
		if isEnumCreate(query) {
			typeName := extractEnumName(query)
			if typeExists(db, typeName) {
				log.Printf("⚡ Enum type %s already exists. Skipping.", typeName)
				continue
			}
		}

		// Execute SQL
		if _, err := db.Exec(query); err != nil {
			log.Printf("❌ ERROR: Failed executing SQL → %v\n    ↳ %s", err, previewQuery(query))
			continue
		}

		logSuccess(query)
	}

	elapsed := time.Since(start)
	green := color.New(color.FgGreen).SprintFunc()
	log.Printf("%s Database migration completed successfully in %s\n", green("✅ DONE:"), elapsed.Round(time.Millisecond))
}

// Check if a query is CREATE TYPE ... AS ENUM
func isEnumCreate(query string) bool {
	matched, _ := regexp.MatchString(`(?i)^CREATE\s+TYPE\s+\w+\s+AS\s+ENUM`, query)
	return matched
}

// Extract enum type name dynamically
func extractEnumName(query string) string {
	re := regexp.MustCompile(`(?i)^CREATE\s+TYPE\s+(\w+)\s+AS\s+ENUM`)
	matches := re.FindStringSubmatch(query)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown_enum"
}

// Check if enum type already exists
func typeExists(db *sql.DB, typeName string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname=$1);`
	if err := db.QueryRow(query, typeName).Scan(&exists); err != nil {
		log.Printf("❌ ERROR checking enum existence %s: %v", typeName, err)
		return false
	}
	return exists
}

// Log success dynamically
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

// Extract table/index/type name dynamically
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

// Short preview for failed query
func previewQuery(query string) string {
	trimmed := strings.ReplaceAll(strings.TrimSpace(query), "\n", " ")
	if len(trimmed) > 120 {
		return trimmed[:120] + "..."
	}
	return trimmed
}
