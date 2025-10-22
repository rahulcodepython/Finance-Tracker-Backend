package models

import (
	"time"

	"github.com/google/uuid"
)

type JwtToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

var JwtTokenColumns = "id, user_id, token, expires_at, created_at"
