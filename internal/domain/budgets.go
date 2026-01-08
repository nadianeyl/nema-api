package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

type Budget struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Name            string
	AvailableBudget decimal.Decimal
	StartDate       time.Time
	EndDate         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int
}

type BudgetItem struct {
	ID          uuid.UUID
	BudgetID    uuid.UUID
	CategoryID  uuid.UUID
	LimitAmount decimal.Decimal
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}
