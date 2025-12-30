package service

import (
	"errors"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type TransactionService struct {
	txProvider *repository.TxProvider
}

func NewTransactionService(txProvider *repository.TxProvider) TransactionService {
	return TransactionService{
		txProvider: txProvider,
	}
}

func (s *TransactionService) Add(req *AddTransactionRequest) (*TransactionResponse, error) {
	var result *domain.Transaction

	err := s.txProvider.WithTx(func(adapters repository.Adapters) error {
		_, err := adapters.CategoryRepo.GetByIDAndTypeForUser(req.CategoryID, req.Type, req.UserID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRecordNotFound):
				return domain.ErrInvalidInputValue
			default:
				return err
			}
		}

		var fromAccount, toAccount *domain.Account

		if req.FromAccountID != nil {
			fromAccount, err = adapters.AccountRepo.GetByIDForUpdate(*req.FromAccountID)
			if err != nil && errors.Is(err, domain.ErrRecordNotFound) {
				return domain.ErrInvalidInputValue
			}

			if err != nil {
				return err
			}

			if fromAccount.UserID != req.UserID {
				return domain.ErrUserNotAllowed
			}
		}

		if req.ToAccountID != nil {
			toAccount, err = adapters.AccountRepo.GetByIDForUpdate(*req.ToAccountID)
			if err != nil && errors.Is(err, domain.ErrRecordNotFound) {
				return domain.ErrInvalidInputValue
			}

			if err != nil {
				return err
			}

			if toAccount.UserID != req.UserID {
				return domain.ErrUserNotAllowed
			}
		}

		transaction := &domain.Transaction{
			UserID:     req.UserID,
			Type:       req.Type,
			CategoryID: req.CategoryID,
			Amount:     req.Amount,
			Date:       req.Date,
		}

		transaction.SetTitle(req.Title)
		transaction.SetNotes(req.Notes)
		transaction.SetFromAccountID(req.FromAccountID)
		transaction.SetToAccountID(req.ToAccountID)

		err = adapters.TransactionRepo.Insert(transaction)
		if err != nil {
			return err
		}

		err = s.updateAccountBalances(adapters, transaction)
		if err != nil {
			return err
		}

		result = transaction
		return nil
	})

	if err != nil {
		return nil, err
	}

	res := &TransactionResponse{
		ID:            result.ID,
		UserID:        result.UserID,
		Type:          result.Type,
		CategoryID:    result.CategoryID,
		Amount:        result.Amount,
		Date:          result.Date,
		Title:         result.GetTitle(),
		Notes:         result.GetNotes(),
		FromAccountID: result.GetFromAccountID(),
		ToAccountID:   result.GetToAccountID(),
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}

	return res, nil
}

func (s *TransactionService) updateAccountBalances(adapters repository.Adapters, transaction *domain.Transaction) error {
	switch transaction.Type {
	case domain.TransactionTypeExpense:
		negativeAmount := transaction.Amount.Neg()
		fromAccountID := transaction.GetFromAccountID()
		return adapters.AccountRepo.UpdateBalance(*fromAccountID, negativeAmount)
	case domain.TransactionTypeIncome:
		toAccountID := transaction.GetToAccountID()
		return adapters.AccountRepo.UpdateBalance(*toAccountID, transaction.Amount)
	case domain.TransactionTypeTransfer:
		negativeAmount := transaction.Amount.Neg()
		fromAccountID := transaction.GetFromAccountID()
		toAccountID := transaction.GetToAccountID()

		err := adapters.AccountRepo.UpdateBalance(*fromAccountID, negativeAmount)
		if err != nil {
			return err
		}

		return adapters.AccountRepo.UpdateBalance(*toAccountID, transaction.Amount)
	default:
		return domain.ErrInvalidTransactionType
	}
}
