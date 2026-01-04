package repository

import (
	"context"
	"database/sql"
)

type db interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repositories struct {
	Users        UserRepository
	Tokens       TokenRepository
	Accounts     AccountRepository
	Categories   CategoryRepository
	Transactions TransactionRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Users:        UserRepository{DB: db},
		Tokens:       TokenRepository{DB: db},
		Accounts:     AccountRepository{DB: db},
		Categories:   CategoryRepository{DB: db},
		Transactions: TransactionRepository{DB: db},
	}
}
