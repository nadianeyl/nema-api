package repository

import (
	"context"
	"time"

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
