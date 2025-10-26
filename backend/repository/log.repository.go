// Package repository provides functionality for interacting with the database.
package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

// CreateLog inserts a new log record into the database.
func CreateLog(log *models.Log, db interfaces.SqlExecutor) error {
	// Construct the SQL query for insertion.
	query := fmt.Sprintf("INSERT INTO logs (%s) VALUES ($1, $2, $3, $4)", models.LogColumns)
	// Execute the query with the log data.
	_, err := db.Exec(query, log.ID, log.UserID, log.Message, log.CreatedAt)
	return err
}

// GetLogsByUserID retrieves all logs for a given user ID from the database, with optional filtering and pagination.
func GetLogsByUserID(userID uuid.UUID, startDate string, endDate string, page int, limit int, db interfaces.SqlExecutor) ([]models.Log, error) {
	// Use a strings.Builder to efficiently construct the SQL query.
	var query strings.Builder
	query.WriteString("SELECT id, user_id, message, created_at FROM logs WHERE user_id = $1")

	// Create a slice to hold the query arguments.
	args := []interface{}{userID}

	// Add date range filtering to the query if both start and end dates are provided.
	if startDate != "" && endDate != "" {
		query.WriteString(" AND created_at BETWEEN $2 AND $3 ORDER BY created_at DESC")
		args = append(args, startDate, endDate)
	}

	// Add pagination to the query.
	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit))

	// Execute the query.
	rows, err := db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan them into Log structs.
	var logs []models.Log
	for rows.Next() {
		var log models.Log
		if err := rows.Scan(&log.ID, &log.UserID, &log.Message, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}