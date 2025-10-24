package database

import (
	"database/sql"
	"log"
	"strings"
)

func CommitExec(query string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Execute within the transaction
	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		log.Println("Transaction rolled back for query", query)
		return err
	}

	// Commit and return
	if err := tx.Commit(); err != nil {
		return err
	}

	if after, ok := strings.CutPrefix(query, "CREATE TABLE IF NOT EXISTS "); ok {
		trimmedQuery := after
		tableName := strings.Split(trimmedQuery, " ")[0]
		log.Printf("Table %s created successfully", tableName)
	}

	return nil
}
