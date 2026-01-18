package service

import (
	"context"
	"time"

	"github.com/govalues/decimal"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type BudgetService struct {
	BudgetRepo      repository.BudgetRepository
	BudgetItemRepo  repository.BudgetItemRepository
	CategoryRepo    repository.CategoryRepository
	TransactionRepo repository.TransactionRepository
}

func NewBudgetService(
	budgetRepo repository.BudgetRepository,
	budgetItemRepo repository.BudgetItemRepository,
	categoryRepo repository.CategoryRepository,
	transactionRepo repository.TransactionRepository,
) BudgetService {
	return BudgetService{
		BudgetRepo:      budgetRepo,
		BudgetItemRepo:  budgetItemRepo,
		CategoryRepo:    categoryRepo,
		TransactionRepo: transactionRepo,
	}
}

func (s *BudgetService) Create(ctx context.Context, req *CreateBudgetRequest) (*BudgetResponse, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	budget := &domain.Budget{
		UserID:          req.UserID,
		Name:            req.Name,
		AvailableBudget: req.AvailableBudget,
		StartDate:       startDate,
		EndDate:         endDate,
	}

	err := s.BudgetRepo.Insert(ctx, budget)
	if err != nil {
		return nil, err
	}

	res := &BudgetResponse{
		ID:              budget.ID,
		UserID:          budget.UserID,
		Name:            budget.Name,
		AvailableBudget: budget.AvailableBudget,
		StartDate:       budget.StartDate,
		EndDate:         budget.EndDate,
		CreatedAt:       budget.CreatedAt,
		UpdatedAt:       budget.UpdatedAt,
	}

	return res, nil
}

func (s *BudgetService) GetByID(ctx context.Context, req *GetBudgetDetailRequest) (*BudgetDetailResponse, error) {
	budget, err := s.BudgetRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if budget.UserID != req.UserID {
		return nil, domain.ErrUserNotAllowed
	}

	items, err := s.BudgetItemRepo.GetByBudgetID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	var itemResponses []BudgetItemDetailResponse
	totalLimitAmount := decimal.Zero
	totalSpentAmount := decimal.Zero

	for _, item := range items {
		category, err := s.CategoryRepo.GetByIDForUser(ctx, item.CategoryID, budget.UserID)
		if err != nil {
			return nil, err
		}

		spentAmount, err := s.TransactionRepo.GetTotalSpentByCategoryAndDateRange(
			ctx,
			item.CategoryID,
			budget.UserID,
			budget.StartDate,
			budget.EndDate,
		)
		if err != nil {
			return nil, err
		}

		remainingAmount, _ := item.LimitAmount.Sub(spentAmount)

		var percentageUsed decimal.Decimal
		if item.LimitAmount.IsPos() {
			quotient, _ := spentAmount.Quo(item.LimitAmount)
			hundred, _ := decimal.New(100, 0)
			percentageUsed, _ = quotient.Mul(hundred)
		}

		itemResponses = append(itemResponses, BudgetItemDetailResponse{
			ID:              item.ID,
			BudgetID:        item.BudgetID,
			CategoryID:      item.CategoryID,
			CategoryName:    category.Name,
			LimitAmount:     item.LimitAmount,
			SpentAmount:     spentAmount,
			RemainingAmount: remainingAmount,
			PercentageUsed:  percentageUsed,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
		})

		totalLimitAmount, _ = totalLimitAmount.Add(item.LimitAmount)
		totalSpentAmount, _ = totalSpentAmount.Add(spentAmount)
	}

	totalRemainingAmount, _ := totalLimitAmount.Sub(totalSpentAmount)

	res := &BudgetDetailResponse{
		ID:                   budget.ID,
		UserID:               budget.UserID,
		Name:                 budget.Name,
		AvailableBudget:      budget.AvailableBudget,
		TotalLimitAmount:     totalLimitAmount,
		TotalSpentAmount:     totalSpentAmount,
		TotalRemainingAmount: totalRemainingAmount,
		StartDate:            budget.StartDate,
		EndDate:              budget.EndDate,
		Items:                itemResponses,
		CreatedAt:            budget.CreatedAt,
		UpdatedAt:            budget.UpdatedAt,
	}

	return res, nil
}
