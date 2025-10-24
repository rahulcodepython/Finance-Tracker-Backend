package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rahulcodepython/finance-tracker-backend/backend/interfaces"
	"github.com/rahulcodepython/finance-tracker-backend/backend/models"
)

func CreateLog(log *models.Log, db interfaces.SqlExecutor) error {
	query := fmt.Sprintf("INSERT INTO logs (%s) VALUES ($1, $2, $3, $4)", models.LogColumns)
	_, err := db.Exec(query, log.ID, log.UserID, log.Message, log.CreatedAt)
	return err
}

func GetLogsByUserID(userID uuid.UUID, startDate string, endDate string, page int, limit int, db interfaces.SqlExecutor) ([]models.Log, error) {
	var query strings.Builder
	query.WriteString("SELECT id, user_id, message, created_at FROM logs WHERE user_id = $1")

	args := []interface{}{userID}
	argCount := 2

	if startDate != "" && endDate != "" {
		query.WriteString(" AND created_at BETWEEN $2 AND $3 ORDER BY created_at DESC")
		args = append(args, "%"+startDate+"%")
		args = append(args, "%"+endDate+"%")
		argCount++
		argCount++
	}

	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit))

	rows, err := db.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
