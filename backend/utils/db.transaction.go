package utils

import (
	"database/sql"
	"log"
)

func DBTransaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Handle rollback in case of panic
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			log.Printf("Transaction panic recovered: %v", r)
		}
	}()

	// Execute your transactional operations
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Commit if everything went fine
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
