package database

import (
	"database/sql"
	"log"
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
		log.Println("Transaction rolled back")
		return err
	}

	// Commit and return
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
