package domain

import "errors"

var (
	ErrDuplicateEmail         = errors.New("duplicate email")
	ErrDuplicateRecord        = errors.New("duplicate record")
	ErrRecordNotFound         = errors.New("record not found")
	ErrEditConflict           = errors.New("edit conflict")
	ErrInvalidInputValue      = errors.New("invalid input value")
	ErrInvalidOrExpiredToken  = errors.New("invalid or expired token")
	ErrInvalidAuthCredentials = errors.New("invalid authentication credentials")
	ErrUserNotAllowed         = errors.New("user not allowed")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)
