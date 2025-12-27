package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID              uuid.UUID
	UserID          uuid.NullUUID
	Name            string
	TransactionType TransactionType
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int
}
