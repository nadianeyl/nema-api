package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	"github.com/nadianeyl/nema-api/internal/domain"
)

var AnonymousUser = &UserResponse{}

type UserResponse struct {
	ID                        uuid.UUID `json:"id"`
	Name                      string    `json:"name"`
	Email                     string    `json:"email"`
	Activated                 bool      `json:"activated"`
	EmailNotificationsEnabled bool      `json:"email_notifications_enabled"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

func (u *UserResponse) IsAnonymous() bool {
	return u == AnonymousUser
}

type TokenResponse struct {
	Plaintext string    `json:"token"`
	Expiry    time.Time `json:"expiry"`
}

type AccountResponse struct {
	ID           uuid.UUID           `json:"id"`
	UserID       uuid.UUID           `json:"user_id"`
	Name         string              `json:"name"`
	Class        domain.AccountClass `json:"class"`
	CurrencyCode string              `json:"currency_code"`
	Balance      decimal.Decimal     `json:"balance"`
	IsBudgeted   bool                `json:"is_budgeted"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type CategoryResponse struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.NullUUID          `json:"user_id"`
	Name            string                 `json:"name"`
	TransactionType domain.TransactionType `json:"transaction_type"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type TransactionResponse struct {
	ID            uuid.UUID              `json:"uuid"`
	UserID        uuid.UUID              `json:"user_id"`
	Type          domain.TransactionType `json:"type"`
	CategoryID    uuid.UUID              `json:"category_id"`
	Amount        decimal.Decimal        `json:"amount"`
	Date          time.Time              `json:"date"`
	Title         *string                `json:"title"`
	Notes         *string                `json:"notes"`
	FromAccountID *uuid.UUID             `json:"from_account_id"`
	ToAccountID   *uuid.UUID             `json:"to_account_id"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}
