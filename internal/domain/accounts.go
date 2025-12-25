package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

type AccountClass string

const (
	AccountClassCCE        AccountClass = "cce"
	AccountClassInvestment AccountClass = "investment"
	AccountClassLiability  AccountClass = "liability"
)

type Account struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Name         string
	Class        AccountClass
	CurrencyCode string
	Balance      decimal.Decimal
	IsBudgeted   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int
}
