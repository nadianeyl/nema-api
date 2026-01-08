package repository

import (
	"context"
	"strings"
	"time"

	"github.com/nadianeyl/nema-api/internal/domain"
)

type BudgetItemRepository struct {
	DB db
}

func (r *BudgetItemRepository) Insert(item *domain.BudgetItem) error {
	query := `
		INSERT INTO budget_items (budget_id, category_id, limit_amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at, version`

	args := []any{item.BudgetID, item.CategoryID, item.LimitAmount}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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
