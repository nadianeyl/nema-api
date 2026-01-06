package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type TransactionRepository struct {
	DB db
}

func (r *TransactionRepository) GetAllForUser(userID uuid.UUID, filters domain.TransactionFilters) ([]*domain.Transaction, domain.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, user_id, type, category_id, amount, date, title, notes, from_account_id, to_account_id, created_at, updated_at, version
		FROM transactions
		WHERE user_id = $1
		AND (type = $2 OR $2 IS NULL)
		AND (category_id = $3 OR $3 IS NULL)
		AND ((from_account_id = $4 OR to_account_id = $4) OR $4 IS NULL)
		AND (date >= $5 OR $5 IS NULL)
		AND (date <= $6 OR $6 IS NULL)
		AND (to_tsvector('simple', COALESCE(title, '')) @@ plainto_tsquery('simple', $7) OR $7 = '')
		ORDER BY %s %s, created_at DESC, id ASC
		LIMIT $8 OFFSET $9`, filters.SortColumn(), filters.SortDirection())

	var typeParam sql.NullString
	if filters.Type != "" {
		typeParam = sql.NullString{String: filters.Type.String(), Valid: true}
	}

	var categoryIDParam sql.NullString
	if filters.CategoryID != "" {
		categoryIDParam = sql.NullString{String: filters.CategoryID, Valid: true}
	}

	var accountIDParam sql.NullString
	if filters.AccountID != "" {
		accountIDParam = sql.NullString{String: filters.AccountID, Valid: true}
	}

	var startDateParam sql.NullString
	if filters.StartDate != "" {
		startDateParam = sql.NullString{String: filters.StartDate, Valid: true}
	}

	var endDateParam sql.NullString
	if filters.EndDate != "" {
		endDateParam = sql.NullString{String: filters.EndDate, Valid: true}
	}

	args := []any{
		userID,
		typeParam,
		categoryIDParam,
		accountIDParam,
		startDateParam,
		endDateParam,
		filters.Title,
		filters.GetLimit(),
		filters.GetOffset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, domain.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	transactions := make([]*domain.Transaction, 0)

	for rows.Next() {
		var transaction domain.Transaction

		err := rows.Scan(
			&totalRecords,
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
			return nil, domain.Metadata{}, err
		}

		transactions = append(transactions, &transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.Metadata{}, err
	}

	metadata := domain.GenerateMetadata(filters.GetLimit(), filters.GetPage(), totalRecords)

	return transactions, metadata, nil
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

func (r *TransactionRepository) GetByIDForUpdate(id uuid.UUID) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, type, category_id, amount, date, title, notes, from_account_id, to_account_id, created_at, updated_at, version
		FROM transactions
		WHERE id = $1
		FOR UPDATE`

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

func (r *TransactionRepository) GetDetailByID(id uuid.UUID) (*domain.TransactionDetail, error) {
	query := `
		SELECT t.id, t.user_id, t.type, t.category_id, t.amount, t.date, t.title, t.notes, t.from_account_id, t.to_account_id, t.created_at, t.updated_at, t.version,
		c.id, c.name, fa.id, fa.name, ta.id, ta.name
		FROM transactions t
		INNER JOIN categories c ON t.category_id = c.id
		LEFT JOIN accounts fa ON t.from_account_id = fa.id
		LEFT JOIN accounts ta ON t.to_account_id = ta.id
		WHERE t.id = $1`

	var (
		detail          domain.TransactionDetail
		fromAccountID   uuid.NullUUID
		fromAccountName sql.NullString
		toAccountID     uuid.NullUUID
		toAccountName   sql.NullString
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		// Transaction fields
		&detail.Transaction.ID,
		&detail.Transaction.UserID,
		&detail.Transaction.Type,
		&detail.Transaction.CategoryID,
		&detail.Transaction.Amount,
		&detail.Transaction.Date,
		&detail.Transaction.Title,
		&detail.Transaction.Notes,
		&detail.Transaction.FromAccountID,
		&detail.Transaction.ToAccountID,
		&detail.Transaction.CreatedAt,
		&detail.Transaction.UpdatedAt,
		&detail.Transaction.Version,
		// Category fields
		&detail.Category.ID,
		&detail.Category.Name,
		// From Account fields (nullable)
		&fromAccountID,
		&fromAccountName,
		// To Account fields (nullable)
		&toAccountID,
		&toAccountName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if fromAccountID.Valid {
		detail.FromAccount = &domain.AccountInfo{
			ID:   fromAccountID.UUID,
			Name: fromAccountName.String,
		}
	}

	if toAccountID.Valid {
		detail.ToAccount = &domain.AccountInfo{
			ID:   toAccountID.UUID,
			Name: toAccountName.String,
		}
	}

	return &detail, nil
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

func (r *TransactionRepository) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM transactions
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
