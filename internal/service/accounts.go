package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type AccountService struct {
	AccountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return AccountService{
		AccountRepo: accountRepo,
	}
}

func (s *AccountService) List(ctx context.Context, req *ListAccountsRequest) ([]*AccountResponse, domain.Metadata, error) {
	accounts, metadata, err := s.AccountRepo.GetAllForUser(ctx, req.UserID, req.Class, req.Filters)
	if err != nil {
		return nil, domain.Metadata{}, err
	}

	res := make([]*AccountResponse, 0)
	for _, account := range accounts {
		res = append(res, &AccountResponse{
			ID:           account.ID,
			UserID:       account.UserID,
			Name:         account.Name,
			Class:        account.Class,
			CurrencyCode: account.CurrencyCode,
			Balance:      account.Balance,
			IsBudgeted:   account.IsBudgeted,
			CreatedAt:    account.CreatedAt,
			UpdatedAt:    account.UpdatedAt,
		})
	}

	return res, metadata, nil
}

func (s *AccountService) Add(ctx context.Context, req *AddAccountRequest) (*AccountResponse, error) {
	account := &domain.Account{
		UserID:       req.UserID,
		Name:         req.Name,
		Class:        req.Class,
		CurrencyCode: req.CurrencyCode,
		Balance:      req.Balance,
		IsBudgeted:   req.IsBudgeted,
	}

	err := s.AccountRepo.Insert(ctx, account)
	if err != nil {
		return nil, err
	}

	res := &AccountResponse{
		ID:           account.ID,
		UserID:       account.UserID,
		Name:         account.Name,
		Class:        account.Class,
		CurrencyCode: account.CurrencyCode,
		Balance:      account.Balance,
		IsBudgeted:   account.IsBudgeted,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}

	return res, nil
}

func (s *AccountService) Update(ctx context.Context, req *UpdateAccountRequest) (*AccountResponse, error) {
	account, err := s.AccountRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.UserID != account.UserID {
		return nil, domain.ErrUserNotAllowed
	}

	if req.Name != nil {
		account.Name = *req.Name
	}

	if req.Class != nil {
		account.Class = *req.Class
	}

	if req.IsBudgeted != nil {
		account.IsBudgeted = *req.IsBudgeted
	}

	err = s.AccountRepo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	res := &AccountResponse{
		ID:           account.ID,
		UserID:       account.UserID,
		Name:         account.Name,
		Class:        account.Class,
		CurrencyCode: account.CurrencyCode,
		Balance:      account.Balance,
		IsBudgeted:   account.IsBudgeted,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}

	return res, nil
}

func (s *AccountService) Delete(ctx context.Context, req *DeleteAccountRequest) error {
	account, err := s.AccountRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if req.UserID != account.UserID {
		return domain.ErrUserNotAllowed
	}

	err = s.AccountRepo.Delete(ctx, account.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountService) GetNetWorth(ctx context.Context, userID uuid.UUID) (*GetNetWorthResponse, error) {
	netWorth, err := s.AccountRepo.GetNetWorthForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := &GetNetWorthResponse{
		CCE:        netWorth.CCE,
		Investment: netWorth.Investment,
		Liability:  netWorth.Liability,
		Total:      netWorth.Total,
	}

	return res, nil
}
