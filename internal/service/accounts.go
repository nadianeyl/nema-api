package service

import (
	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type AccountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return AccountService{
		accountRepo: accountRepo,
	}
}

func (s *AccountService) List(req *ListAccountsRequest) ([]*AccountResponse, domain.Metadata, error) {
	accounts, metadata, err := s.accountRepo.GetAllForUser(req.UserID, req.Class, req.Filters)
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

func (s *AccountService) Add(req *AddAccountRequest) (*AccountResponse, error) {
	account := &domain.Account{
		UserID:       req.UserID,
		Name:         req.Name,
		Class:        req.Class,
		CurrencyCode: req.CurrencyCode,
		Balance:      req.Balance,
		IsBudgeted:   req.IsBudgeted,
	}

	err := s.accountRepo.Insert(account)
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
