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

type AccountInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CategoryResponse struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.NullUUID          `json:"user_id"`
	Name            string                 `json:"name"`
	TransactionType domain.TransactionType `json:"transaction_type"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type CategoryInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
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

type TransactionDetailResponse struct {
	Transaction TransactionResponse `json:"transaction"`
	Category    CategoryInfo        `json:"category"`
	FromAccount *AccountInfo        `json:"from_account"`
	ToAccount   *AccountInfo        `json:"to_account"`
}

type BudgetResponse struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	Name            string          `json:"name"`
	AvailableBudget decimal.Decimal `json:"available_budget"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type BudgetItemResponse struct {
	ID          uuid.UUID       `json:"id"`
	BudgetID    uuid.UUID       `json:"budget_id"`
	CategoryID  uuid.UUID       `json:"category_id"`
	LimitAmount decimal.Decimal `json:"limit_amount"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type BudgetDetailResponse struct {
	ID                   uuid.UUID                  `json:"id"`
	UserID               uuid.UUID                  `json:"user_id"`
	Name                 string                     `json:"name"`
	AvailableBudget      decimal.Decimal            `json:"available_budget"`
	TotalLimitAmount     decimal.Decimal            `json:"total_limit_amount"`
	TotalSpentAmount     decimal.Decimal            `json:"total_spent_amount"`
	TotalRemainingAmount decimal.Decimal            `json:"total_remaining_amount"`
	StartDate            time.Time                  `json:"start_date"`
	EndDate              time.Time                  `json:"end_date"`
	Items                []BudgetItemDetailResponse `json:"items"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
}

type BudgetItemDetailResponse struct {
	ID              uuid.UUID       `json:"id"`
	BudgetID        uuid.UUID       `json:"budget_id"`
	CategoryID      uuid.UUID       `json:"category_id"`
	CategoryName    string          `json:"category_name"`
	LimitAmount     decimal.Decimal `json:"limit_amount"`
	SpentAmount     decimal.Decimal `json:"spent_amount"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	PercentageUsed  decimal.Decimal `json:"percentage_used"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type GetNetWorthResponse struct {
	CCE        decimal.Decimal `json:"cce"`
	Investment decimal.Decimal `json:"investment"`
	Liability  decimal.Decimal `json:"liability"`
	Total      decimal.Decimal `json:"total"`
}
