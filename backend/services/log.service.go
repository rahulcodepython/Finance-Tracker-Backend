package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
	"github.com/rahulcodepython/finance-tracker-backend/backend/repository"
	"github.com/rahulcodepython/finance-tracker-backend/backend/utils"
)

func CreateLog(userID uuid.UUID, message string, db *sql.DB) error {
	log := &models.Log{
		ID:        uuid.New(),
		UserID:    userID,
		Message:   message,
		CreatedAt: time.Now().In(utils.LOC),
	}
	return repository.CreateLog(log, db)
}

func GetLogs(userID uuid.UUID, startDate string, endDate string, page int, limit int, db *sql.DB) ([]models.Log, error) {
	return repository.GetLogsByUserID(userID, startDate, endDate, page, limit, db)
}
