package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

const LogColumns = "id, user_id, message, created_at"
