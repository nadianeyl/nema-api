package repository

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type BudgetItemRepository struct {
	DB db
}

func (r *BudgetItemRepository) Insert(ctx context.Context, item *domain.BudgetItem) error {
	query := `
		INSERT INTO budget_items (budget_id, category_id, limit_amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at, version`

	args := []any{item.BudgetID, item.CategoryID, item.LimitAmount}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.Version,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "uc_budget_items_budget_category"):
			return domain.ErrDuplicateRecord
		default:
			return err
		}
	}

	return nil
}

func (r *BudgetItemRepository) GetByBudgetID(ctx context.Context, budgetID uuid.UUID) ([]*domain.BudgetItem, error) {
	query := `
		SELECT id, budget_id, category_id, limit_amount, created_at, updated_at, version
		FROM budget_items
		WHERE budget_id = $1
		ORDER BY created_at ASC, id ASC`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, query, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.BudgetItem

	for rows.Next() {
		var item domain.BudgetItem

		err := rows.Scan(
			&item.ID,
			&item.BudgetID,
			&item.CategoryID,
			&item.LimitAmount,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Version,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
