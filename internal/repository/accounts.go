package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type AccountRepository struct {
	DB *sql.DB
}

func (r *AccountRepository) Insert(account *domain.Account) error {
	query := `
		INSERT INTO accounts (user_id, name, class, currency_code, balance, is_budgeted)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at, version`

	args := []any{
		account.UserID,
		account.Name,
		account.Class,
		account.CurrencyCode,
		account.Balance,
		account.IsBudgeted,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&account.ID,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Version,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "pq: invalid input value for enum"):
			return domain.ErrInvalidInputValue
		default:
			return err
		}
	}

	return nil
}
