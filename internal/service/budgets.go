package service

import (
	"time"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type BudgetService struct {
	BudgetRepo repository.BudgetRepository
}

func NewBudgetService(budgetRepo repository.BudgetRepository) BudgetService {
	return BudgetService{
		BudgetRepo: budgetRepo,
	}
}

func (s *BudgetService) Create(req *CreateBudgetRequest) (*BudgetResponse, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	budget := &domain.Budget{
		UserID:          req.UserID,
		Name:            req.Name,
		AvailableBudget: req.AvailableBudget,
		StartDate:       startDate,
		EndDate:         endDate,
	}

	err := s.BudgetRepo.Insert(budget)
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
