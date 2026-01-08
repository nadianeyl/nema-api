package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/validator"
)

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "email is required")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "email is invalid")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "password is required")
	v.Check(len(password) >= 8, "password", "password should be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "password should not be more than 72 bytes long")
}

func ValidateRegisterUserReq(v *validator.Validator, req *RegisterUserRequest) {
	v.Check(len(req.Name) <= 500, "name", "name should not be more than 500 bytes long")

	ValidateEmail(v, req.Email)

	ValidatePasswordPlaintext(v, req.Password)
}

type ActivateUserRequest struct {
	TokenPlaintext string `json:"token"`
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "token is required")
	v.Check(len(tokenPlaintext) == 26, "token", "token must be 26 bytes long")
}

type CreateAuthTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ListAccountsRequest struct {
	UserID uuid.UUID           `json:"-"`
	Class  domain.AccountClass `json:"class"`
	domain.Filters
}

func ValidateListAccountsReq(v *validator.Validator, req *ListAccountsRequest) {
	if req.Class != "" {
		v.Check(validator.PermittedValue(req.Class, domain.GetAccountClasses()...), "class", "class is invalid")
	}

	ValidateFilters(v, req.Filters)
}

type AddAccountRequest struct {
	UserID       uuid.UUID           `json:"-"`
	Name         string              `json:"name"`
	Class        domain.AccountClass `json:"class"`
	CurrencyCode string              `json:"currency_code"`
	Balance      decimal.Decimal     `json:"balance"`
	IsBudgeted   bool                `json:"is_budgeted"`
}

func ValidateAddAccountReq(v *validator.Validator, req *AddAccountRequest) {
	v.Check(req.Name != "", "name", "name is required")

	v.Check(req.Class != "", "class", "class is required")
	v.Check(validator.PermittedValue(req.Class, domain.GetAccountClasses()...), "class", "class is invalid")

	v.Check(req.CurrencyCode != "", "currency_code", "currency code is required")
	v.Check(len(req.CurrencyCode) == 3, "currency_code", "currency code must be 3 characters long")
	v.Check(validator.PermittedValue(req.CurrencyCode, domain.GetCurrencyCodes()...), "currency_code", "currency code is invalid")
}

type UpdateAccountRequest struct {
	ID         uuid.UUID            `json:"-"`
	UserID     uuid.UUID            `json:"-"`
	Name       *string              `json:"name"`
	Class      *domain.AccountClass `json:"class"`
	IsBudgeted *bool                `json:"is_budgeted"`
}

func ValidateUpdateAccountReq(v *validator.Validator, req *UpdateAccountRequest) {
	if req.Name != nil {
		v.Check(*req.Name != "", "name", "name is required")
	}

	if req.Class != nil {
		v.Check(*req.Class != "", "class", "class is required")
		v.Check(validator.PermittedValue(*req.Class, domain.GetAccountClasses()...), "class", "class is invalid")
	}
}

type DeleteAccountRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
}

func ValidateFilters(v *validator.Validator, f domain.Filters) {
	v.Check(f.Limit > 0, "limit", "limit must be greater than zero")
	v.Check(f.Limit <= 100, "limit", "limit must be a maximum of 100")
	v.Check(f.Page > 0, "page", "page must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "page must be a maximum of 10 million")
}

type ListCategoriesRequest struct {
	UserID          uuid.UUID              `json:"-"`
	TransactionType domain.TransactionType `json:"transaction_type"`
	domain.Filters
}

func ValidateListCategoriesReq(v *validator.Validator, req *ListCategoriesRequest) {
	if req.TransactionType != "" {
		v.Check(validator.PermittedValue(req.TransactionType, domain.GetTransactionTypes()...), "transaction_type", "transaction type is invalid")
	}

	ValidateFilters(v, req.Filters)
}

type ListTransactionsRequest struct {
	UserID uuid.UUID `json:"-"`
	domain.TransactionFilters
}

func ValidateListTransactionsReq(v *validator.Validator, req *ListTransactionsRequest) {
	if req.Type != "" {
		v.Check(validator.PermittedValue(req.Type, domain.GetTransactionTypes()...), "type", "type is invalid")
	}

	if req.CategoryID != "" {
		_, err := uuid.Parse(req.CategoryID)
		v.Check(err == nil, "category_id", "category ID must be a valid UUID")
	}

	if req.AccountID != "" {
		_, err := uuid.Parse(req.AccountID)
		v.Check(err == nil, "account_id", "account ID must be a valid UUID")
	}

	if req.StartDate != "" {
		_, err := time.Parse(time.RFC3339, req.StartDate)
		v.Check(err == nil, "start_date", "start date must be a valid RFC3339 timestamp")
	}

	if req.EndDate != "" {
		_, err := time.Parse(time.RFC3339, req.EndDate)
		v.Check(err == nil, "end_date", "end date must be a valid RFC3339 timestamp")
	}

	if req.StartDate != "" && req.EndDate != "" {
		startDate, _ := time.Parse(time.RFC3339, req.StartDate)
		endDate, _ := time.Parse(time.RFC3339, req.EndDate)
		v.Check(!startDate.After(endDate), "end_date", "end date must be after start date")
	}

	if req.Title != "" {
		v.Check(len(req.Title) <= 100, "title", "title must not be more than 100 characters")
	}

	ValidateFilters(v, req.Filters)
}

type AddTransactionRequest struct {
	UserID        uuid.UUID              `json:"-"`
	Type          domain.TransactionType `json:"type"`
	CategoryID    uuid.UUID              `json:"category_id"`
	Amount        decimal.Decimal        `json:"amount"`
	Date          time.Time              `json:"date"`
	Title         *string                `json:"title"`
	Notes         *string                `json:"notes"`
	FromAccountID *uuid.UUID             `json:"from_account_id"`
	ToAccountID   *uuid.UUID             `json:"to_account_id"`
}

func ValidateAccountID(v *validator.Validator, transactionType domain.TransactionType, fromAccountID, toAccountID *uuid.UUID) {
	if transactionType == domain.TransactionTypeExpense {
		v.Check(fromAccountID != nil, "from_account_id", "source account ID is required")
	}

	if transactionType == domain.TransactionTypeIncome {
		v.Check(toAccountID != nil, "to_account_id", "destination account ID is required")
	}

	if transactionType == domain.TransactionTypeTransfer {
		v.Check(fromAccountID != nil, "from_account_id", "source account ID is required")
		v.Check(toAccountID != nil, "to_account_id", "destination account ID is required")

		if fromAccountID != nil && toAccountID != nil {
			v.Check(*fromAccountID != *toAccountID, "to_account_id", "destination account ID must be different from source account ID")
		}
	}
}

func ValidateAddTransactionReq(v *validator.Validator, req *AddTransactionRequest) {
	v.Check(req.Type != "", "type", "type is required")
	v.Check(validator.PermittedValue(req.Type, domain.GetTransactionTypes()...), "type", "type is invalid")

	v.Check(req.CategoryID != uuid.Nil, "category_id", "category ID is required")

	v.Check(req.Amount.IsPos(), "amount", "amount must be greater than zero")

	v.Check(!req.Date.IsZero(), "date", "date is required")

	if req.Title != nil {
		v.Check(len(*req.Title) <= 100, "title", "title must not be more than 100 characters")
	}

	if req.Notes != nil {
		v.Check(len(*req.Notes) <= 400, "notes", "notes must not be more than 400 characters")
	}

	ValidateAccountID(v, req.Type, req.FromAccountID, req.ToAccountID)
}

type GetTransactionDetailRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
}

type UpdateTransactionRequest struct {
	ID            uuid.UUID               `json:"-"`
	UserID        uuid.UUID               `json:"-"`
	Type          *domain.TransactionType `json:"type"`
	CategoryID    *uuid.UUID              `json:"category_id"`
	Amount        *decimal.Decimal        `json:"amount"`
	Date          *time.Time              `json:"date"`
	Title         *string                 `json:"title"`
	Notes         *string                 `json:"notes"`
	FromAccountID *uuid.UUID              `json:"from_account_id"`
	ToAccountID   *uuid.UUID              `json:"to_account_id"`
}

func ValidateUpdateTransactionReq(v *validator.Validator, req *UpdateTransactionRequest) {
	if req.Type != nil {
		v.Check(*req.Type != "", "type", "type is required")
		v.Check(validator.PermittedValue(*req.Type, domain.GetTransactionTypes()...), "type", "type is invalid")
	}

	if req.CategoryID != nil {
		v.Check(*req.CategoryID != uuid.Nil, "category_id", "category ID is required")
	}

	if req.Amount != nil {
		v.Check(req.Amount.IsPos(), "amount", "amount must be greater than zero")
	}

	if req.Date != nil {
		v.Check(!req.Date.IsZero(), "date", "date is required")
	}

	if req.Title != nil {
		v.Check(len(*req.Title) <= 100, "title", "title must not be more than 100 characters")
	}

	if req.Notes != nil {
		v.Check(len(*req.Notes) <= 400, "notes", "notes must not be more than 400 characters")
	}
}

func ApplyTransactionUpdates(transaction *domain.Transaction, req *UpdateTransactionRequest) {
	if req.Type != nil {
		transaction.Type = *req.Type
	}

	if req.CategoryID != nil {
		transaction.CategoryID = *req.CategoryID
	}

	if req.Amount != nil {
		transaction.Amount = *req.Amount
	}

	if req.Date != nil {
		transaction.Date = *req.Date
	}

	if req.Title != nil {
		transaction.SetTitle(req.Title)
	}

	if req.Notes != nil {
		transaction.SetNotes(req.Notes)
	}

	if req.FromAccountID != nil {
		transaction.SetFromAccountID(req.FromAccountID)
	}

	if req.ToAccountID != nil {
		transaction.SetToAccountID(req.ToAccountID)
	}
}

type DeleteTransactionRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
}

type CreateBudgetRequest struct {
	UserID          uuid.UUID       `json:"-"`
	Name            string          `json:"name"`
	AvailableBudget decimal.Decimal `json:"available_budget"`
	StartDate       string          `json:"start_date"` // YYYY-MM-DD format
	EndDate         string          `json:"end_date"`   // YYYY-MM-DD format
}

func ValidateCreateBudgetReq(v *validator.Validator, req *CreateBudgetRequest) {
	v.Check(req.Name != "", "name", "name is required")
	v.Check(len(req.Name) <= 100, "name", "name must not be more than 100 characters")

	v.Check(!req.AvailableBudget.IsNeg(), "available_budget", "available budget must be zero or positive")

	v.Check(req.StartDate != "", "start_date", "start date is required")
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	v.Check(err == nil, "start_date", "start date must be in YYYY-MM-DD format")

	v.Check(req.EndDate != "", "end_date", "end date is required")
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	v.Check(err == nil, "end_date", "end date must be in YYYY-MM-DD format")

	if err == nil {
		v.Check(endDate.After(startDate), "end_date", "end date must be after start date")
	}
}

type CreateBudgetItemRequest struct {
	BudgetID    uuid.UUID       `json:"-"`
	UserID      uuid.UUID       `json:"-"`
	CategoryID  uuid.UUID       `json:"category_id"`
	LimitAmount decimal.Decimal `json:"limit_amount"`
}

func ValidateCreateBudgetItemReq(v *validator.Validator, req *CreateBudgetItemRequest) {
	v.Check(req.CategoryID != uuid.Nil, "category_id", "category ID is required")

	v.Check(req.LimitAmount.IsPos(), "limit_amount", "limit amount must be positive")
}

type GetBudgetDetailRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID uuid.UUID `json:"-"`
}
