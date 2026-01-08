package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type BudgetRepository struct {
	DB db
}

func (r *BudgetRepository) Insert(budget *domain.Budget) error {
	query := `
		INSERT INTO budgets (user_id, name, available_budget, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, version`

	args := []any{
		budget.UserID,
		budget.Name,
		budget.AvailableBudget,
		budget.StartDate,
		budget.EndDate,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&budget.ID,
		&budget.CreatedAt,
		&budget.UpdatedAt,
		&budget.Version,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "uc_budgets_user_date_range"):
			return domain.ErrDuplicateRecord
		default:
			return err
		}
	}

	return nil
}

func (r *BudgetRepository) GetByID(id uuid.UUID) (*domain.Budget, error) {
	query := `
		SELECT id, user_id, name, available_budget, start_date, end_date, created_at, updated_at, version
		FROM budgets
		WHERE id = $1`

	var budget domain.Budget

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.Name,
		&budget.AvailableBudget,
		&budget.StartDate,
		&budget.EndDate,
		&budget.CreatedAt,
		&budget.UpdatedAt,
		&budget.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &budget, nil
}
