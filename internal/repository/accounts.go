package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type AccountRepository struct {
	DB db
}

func (r *AccountRepository) GetAllForUser(ctx context.Context, userID uuid.UUID, class domain.AccountClass, filters domain.Filters) ([]*domain.Account, domain.Metadata, error) {
	query := `
		SELECT COUNT(*) OVER(), id, user_id, name, class, currency_code, balance, is_budgeted, created_at, updated_at, version
		FROM accounts
		WHERE user_id = $1 
		AND (class = $2 OR $2 IS NULL)
		ORDER BY created_at DESC, id ASC
		LIMIT $3 OFFSET $4`

	var classParam sql.NullString
	if class != "" {
		classParam = sql.NullString{String: class.String(), Valid: true}
	}

	args := []any{userID, classParam, filters.GetLimit(), filters.GetOffset()}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
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

func (r *AccountRepository) Insert(ctx context.Context, account *domain.Account) error {
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

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
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

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	query := `
		SELECT id, user_id, name, class, currency_code, balance, is_budgeted, created_at, updated_at, version
		FROM accounts
		WHERE id = $1`

	var account domain.Account

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &account, nil
}

func (r *AccountRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	query := `
		SELECT id, user_id, name, class, currency_code, balance, is_budgeted, created_at, updated_at, version
		FROM accounts
		WHERE id = $1
		FOR UPDATE`

	var account domain.Account

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &account, nil
}

func (r *AccountRepository) Update(ctx context.Context, account *domain.Account) error {
	query := `
		UPDATE accounts
		SET name = $1, class = $2, is_budgeted = $3, updated_at = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING updated_at, version`

	args := []any{
		account.Name,
		account.Class,
		account.IsBudgeted,
		time.Now(),
		account.ID,
		account.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&account.UpdatedAt, &account.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amountChange decimal.Decimal) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1, updated_at = $2
		WHERE id = $3`

	args := []any{amountChange, time.Now(), id}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrRecordNotFound
	}

	return nil
}

func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM accounts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrRecordNotFound
	}

	return nil
}

func (r *AccountRepository) GetNetWorthForUser(ctx context.Context, userID uuid.UUID) (*domain.NetWorth, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN class = 'cce' THEN balance ELSE 0 END), 0) as total_cce,
			COALESCE(SUM(CASE WHEN class = 'investment' THEN balance ELSE 0 END), 0) as total_investment,
			COALESCE(SUM(CASE WHEN class = 'liability' THEN balance ELSE 0 END), 0) as total_liability
		FROM accounts
		WHERE user_id = $1`

	var netWorth domain.NetWorth

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&netWorth.CCE,
		&netWorth.Investment,
		&netWorth.Liability,
	)
	if err != nil {
		return nil, err
	}

	assets, _ := netWorth.CCE.Add(netWorth.Investment)
	netWorth.Total, _ = assets.Sub(netWorth.Liability)

	return &netWorth, nil
}
