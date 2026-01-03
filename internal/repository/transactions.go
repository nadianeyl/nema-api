package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type TransactionRepository struct {
	DB db
}

func (r *TransactionRepository) Insert(transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, type, category_id, amount, date, title, notes, from_account_id, to_account_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at, version`

	args := []any{
		transaction.UserID,
		transaction.Type,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Date,
		transaction.Title,
		transaction.Notes,
		transaction.FromAccountID,
		transaction.ToAccountID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) GetByID(id uuid.UUID) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, type, category_id, amount, date, title, notes, from_account_id, to_account_id, created_at, updated_at, version
		FROM transactions
		WHERE id = $1`

	var transaction domain.Transaction

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Type,
		&transaction.CategoryID,
		&transaction.Amount,
		&transaction.Date,
		&transaction.Title,
		&transaction.Notes,
		&transaction.FromAccountID,
		&transaction.ToAccountID,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &transaction, nil
}

func (r *TransactionRepository) Update(transaction *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET type = $1, category_id = $2, amount = $3, date = $4, title = $5, notes = $6, from_account_id = $7, to_account_id = $8, updated_at = $9, version = version + 1
		WHERE id = $10 AND version = $11
		RETURNING updated_at, version`

	args := []any{
		transaction.Type,
		transaction.CategoryID,
		transaction.Amount,
		transaction.Date,
		transaction.Title,
		transaction.Notes,
		transaction.FromAccountID,
		transaction.ToAccountID,
		time.Now(),
		transaction.ID,
		transaction.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&transaction.UpdatedAt, &transaction.Version)
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
