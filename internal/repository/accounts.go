package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type AccountRepository struct {
	DB *sql.DB
}

func (r *AccountRepository) GetAllForUser(userID uuid.UUID, class domain.AccountClass, filters domain.Filters) ([]*domain.Account, domain.Metadata, error) {
	query := `
		SELECT COUNT(*) OVER(), id, user_id, name, class, currency_code, balance, is_budgeted, created_at, updated_at, version
		FROM accounts
		WHERE user_id = $1 
		AND class = $2 OR $2 IS NULL
		ORDER BY created_at, id ASC
		LIMIT $3 OFFSET $4`

	var classParam sql.NullString
	if class != "" {
		classParam = sql.NullString{String: class.String(), Valid: true}
	}

	args := []any{userID, classParam, filters.GetLimit(), filters.GetOffset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, domain.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	accounts := make([]*domain.Account, 0)

	for rows.Next() {
		var account domain.Account

		err := rows.Scan(
			&totalRecords,
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.Class,
			&account.CurrencyCode,
			&account.Balance,
			&account.IsBudgeted,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.Version,
		)
		if err != nil {
			return nil, domain.Metadata{}, err
		}

		accounts = append(accounts, &account)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.Metadata{}, err
	}

	metadata := domain.GenerateMetadata(filters.GetLimit(), filters.GetPage(), totalRecords)

	return accounts, metadata, nil
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
