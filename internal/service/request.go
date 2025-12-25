package service

import (
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

func ValidateFilters(v *validator.Validator, f domain.Filters) {
	v.Check(f.Limit > 0, "limit", "limit must be greater than zero")
	v.Check(f.Limit <= 100, "limit", "limit must be a maximum of 100")
	v.Check(f.Page > 0, "page", "page must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "page must be a maximum of 10 million")
}
