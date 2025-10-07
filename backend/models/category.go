package models

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID       `json:"id"`
	Name string          `json:"name"`
	Type TransactionType `json:"type"`
}

var CategoryColumns = "id, name, type"
