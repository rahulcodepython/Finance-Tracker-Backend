// Package services provides business logic for the application.
package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

// CreateLog creates a new log entry for a user.
func CreateLog(userID uuid.UUID, message string, db *sql.DB) error {
	// Create a new Log model.
	log := &models.Log{
		ID:        uuid.New(),
		UserID:    userID,
		Message:   message,
		CreatedAt: time.Now().In(utils.LOC),
	}
	// Create the log entry in the database.
	return repository.CreateLog(log, db)
}

// GetLogs retrieves all logs for a user with optional filtering and pagination.
func GetLogs(userID uuid.UUID, startDate string, endDate string, page int, limit int, db *sql.DB) ([]models.Log, error) {
	return repository.GetLogsByUserID(userID, startDate, endDate, page, limit, db)
}