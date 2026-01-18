package service

import (
	"context"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type BudgetItemService struct {
	BudgetItemRepo repository.BudgetItemRepository
	BudgetRepo     repository.BudgetRepository
	CategoryRepo   repository.CategoryRepository
}

func NewBudgetItemService(budgetItemRepo repository.BudgetItemRepository, budgetRepo repository.BudgetRepository, categoryRepo repository.CategoryRepository) BudgetItemService {
	return BudgetItemService{
		BudgetItemRepo: budgetItemRepo,
		BudgetRepo:     budgetRepo,
		CategoryRepo:   categoryRepo,
	}
}

func (s *BudgetItemService) Create(ctx context.Context, req *CreateBudgetItemRequest) (*BudgetItemResponse, error) {
	budget, err := s.BudgetRepo.GetByID(ctx, req.BudgetID)
	if err != nil {
		return nil, err
	}

	if budget.UserID != req.UserID {
		return nil, domain.ErrUserNotAllowed
	}

	_, err = s.CategoryRepo.GetByIDAndTypeForUser(ctx, req.CategoryID, domain.TransactionTypeExpense, budget.UserID)
	if err != nil {
		return nil, err
	}

	item := &domain.BudgetItem{
		BudgetID:    req.BudgetID,
		CategoryID:  req.CategoryID,
		LimitAmount: req.LimitAmount,
	}

	err = s.BudgetItemRepo.Insert(ctx, item)
	if err != nil {
		return nil, err
	}

	res := &BudgetItemResponse{
		ID:          item.ID,
		BudgetID:    item.BudgetID,
		CategoryID:  item.CategoryID,
		LimitAmount: item.LimitAmount,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	return res, nil
}
