package domain

import "errors"

var (
	ErrDuplicateEmail         = errors.New("duplicate email")
	ErrRecordNotFound         = errors.New("record not found")
	ErrEditConflict           = errors.New("edit conflict")
	ErrInvalidOrExpiredToken  = errors.New("invalid or expired token")
	ErrInvalidAuthCredentials = errors.New("invalid authentication credentials")
)
