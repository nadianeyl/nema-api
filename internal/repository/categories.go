package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type CategoryRepository struct {
	DB db
}

func (r *CategoryRepository) GetAllForUser(ctx context.Context, userID uuid.UUID, transactionType domain.TransactionType, filters domain.Filters) ([]*domain.Category, domain.Metadata, error) {
	query := `
		SELECT COUNT(*) OVER(), id, user_id, name, transaction_type, created_at, updated_at, version
		FROM categories
		WHERE (user_id = $1 OR user_id IS NULL)
		AND (transaction_type = $2 OR $2 IS NULL)
		ORDER BY transaction_type ASC, name ASC, id ASC
		LIMIT $3 OFFSET $4`

	var transactionTypeParam sql.NullString
	if transactionType != "" {
		transactionTypeParam = sql.NullString{String: transactionType.String(), Valid: true}
	}

	args := []any{userID, transactionTypeParam, filters.GetLimit(), filters.GetOffset()}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, domain.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	categories := make([]*domain.Category, 0)

	for rows.Next() {
		var category domain.Category

		err = rows.Scan(
			&totalRecords,
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.TransactionType,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.Version,
		)
		if err != nil {
			return nil, domain.Metadata{}, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.Metadata{}, err
	}

	metadata := domain.GenerateMetadata(filters.GetLimit(), filters.GetPage(), totalRecords)

	return categories, metadata, nil
}

func (r *CategoryRepository) GetByIDForUser(ctx context.Context, id, userID uuid.UUID) (*domain.Category, error) {
	query := `
		SELECT id, user_id, name, transaction_type, created_at, updated_at, version
		FROM categories
		WHERE id = $1
		AND (user_id = $2 OR user_id IS NULL)`

	var category domain.Category

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.TransactionType,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (r *CategoryRepository) GetByIDAndTypeForUser(ctx context.Context, id uuid.UUID, transactionType domain.TransactionType, userID uuid.UUID) (*domain.Category, error) {
	query := `
		SELECT id, user_id, name, transaction_type, created_at, updated_at, version
		FROM categories
		WHERE id = $1
		AND transaction_type = $2
		AND (user_id = $3 OR user_id IS NULL)
	`

	args := []any{id, transactionType, userID}

	var category domain.Category

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.TransactionType,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}
