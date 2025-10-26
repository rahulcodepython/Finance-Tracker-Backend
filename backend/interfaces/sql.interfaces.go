// Package interfaces defines shared interfaces for the application.
package interfaces

import "database/sql"

// SqlExecutor defines an interface for executing SQL queries.
// This allows for easy mocking of database operations in tests.
type SqlExecutor interface {
	// Exec executes a query without returning any rows.
	Exec(query string, args ...interface{}) (sql.Result, error)
	// QueryRow executes a query that is expected to return at most one row.
	QueryRow(query string, args ...interface{}) *sql.Row
	// Query executes a query that returns rows, typically a SELECT.
	Query(query string, args ...interface{}) (*sql.Rows, error)
}