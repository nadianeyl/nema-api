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

func (c AccountClass) String() string {
	return string(c)
}

func GetAccountClasses() []AccountClass {
	return []AccountClass{AccountClassCCE, AccountClassInvestment, AccountClassLiability}
}

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

type AccountInfo struct {
	ID   uuid.UUID
	Name string
}

type NetWorth struct {
	CCE        decimal.Decimal
	Investment decimal.Decimal
	Liability  decimal.Decimal
	Total      decimal.Decimal
}
