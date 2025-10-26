// Package utils provides utility functions for the application.
package utils

import (
	"database/sql"
	"log"
)

// DBTransaction is a utility function that wraps a database transaction.
// It begins a transaction, executes a function within that transaction, and then commits or rolls back the transaction based on the outcome of the function.
func DBTransaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	// Begin a new database transaction.
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Use a deferred function to handle panics and roll back the transaction if one occurs.
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			log.Printf("Transaction panic recovered: %v", r)
		}
	}()

	// Execute the provided function within the transaction.
	if err := fn(tx); err != nil {
		// If the function returns an error, roll back the transaction.
		_ = tx.Rollback()
		return err
	}

	// If the function completes without error, commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}