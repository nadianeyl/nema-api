package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
	"github.com/nadianeyl/nema-api/internal/validator"
)

type TransactionService struct {
	TxProvider *repository.TxProvider
}

func NewTransactionService(txProvider *repository.TxProvider) TransactionService {
	return TransactionService{
		TxProvider: txProvider,
	}
}

func (s *TransactionService) Add(req *AddTransactionRequest) (*TransactionResponse, error) {
	var result *domain.Transaction

	err := s.TxProvider.WithTx(func(adapters repository.Adapters) error {
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

func (s *TransactionService) Update(req *UpdateTransactionRequest, v *validator.Validator) (*TransactionResponse, error) {
	var result *domain.Transaction

	err := s.TxProvider.WithTx(func(adapters repository.Adapters) error {
		transaction, err := adapters.TransactionRepo.GetByID(req.ID)
		if err != nil {
			return err
		}

		if transaction.UserID != req.UserID {
			return domain.ErrUserNotAllowed
		}

		oldType := transaction.Type
		oldAmount := transaction.Amount
		oldFromAccountID := transaction.GetFromAccountID()
		oldToAccountID := transaction.GetToAccountID()

		ApplyTransactionUpdates(transaction, req)

		newFromAccountID := transaction.GetFromAccountID()
		newToAccountID := transaction.GetToAccountID()

		if ValidateAccountID(v, transaction.Type, newFromAccountID, newToAccountID); !v.Valid() {
			return domain.ErrInvalidInputValue
		}

		_, err = adapters.CategoryRepo.GetByIDAndTypeForUser(transaction.CategoryID, transaction.Type, req.UserID)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRecordNotFound):
				v.AddError("category_id", "category ID is invalid")
				return domain.ErrInvalidInputValue
			default:
				return err
			}
		}

		accountIDs := make(map[uuid.UUID]bool)
		if oldFromAccountID != nil {
			accountIDs[*oldFromAccountID] = true
		}
		if oldToAccountID != nil {
			accountIDs[*oldToAccountID] = true
		}
		if newFromAccountID != nil {
			accountIDs[*newFromAccountID] = true
		}
		if newToAccountID != nil {
			accountIDs[*newToAccountID] = true
		}

		for accountID := range accountIDs {
			account, err := adapters.AccountRepo.GetByIDForUpdate(accountID)
			if err != nil && errors.Is(err, domain.ErrRecordNotFound) {
				v.AddError("account_id", "source account ID or destination account ID is invalid")
				return domain.ErrInvalidInputValue
			}

			if err != nil {
				return err
			}

			if account.UserID != req.UserID {
				return domain.ErrUserNotAllowed
			}
		}

		err = s.revertAccountBalances(adapters, oldType, oldAmount, oldFromAccountID, oldToAccountID)
		if err != nil {
			return err
		}

		switch transaction.Type {
		case domain.TransactionTypeExpense:
			transaction.SetToAccountID(nil)
		case domain.TransactionTypeIncome:
			transaction.SetFromAccountID(nil)
		}

		err = adapters.TransactionRepo.Update(transaction)
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

func (s *TransactionService) Delete(req *DeleteTransactionRequest) error {
	err := s.TxProvider.WithTx(func(adapters repository.Adapters) error {
		transaction, err := adapters.TransactionRepo.GetByIDForUpdate(req.ID)
		if err != nil {
			return err
		}

		if transaction.UserID != req.UserID {
			return domain.ErrUserNotAllowed
		}

		err = s.revertAccountBalances(adapters, transaction.Type, transaction.Amount, transaction.GetFromAccountID(), transaction.GetToAccountID())
		if err != nil {
			return err
		}

		err = adapters.TransactionRepo.Delete(transaction.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
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

func (s *TransactionService) revertAccountBalances(adapters repository.Adapters, transactionType domain.TransactionType, amount decimal.Decimal, fromAccountID, toAccountID *uuid.UUID) error {
	switch transactionType {
	case domain.TransactionTypeExpense:
		return adapters.AccountRepo.UpdateBalance(*fromAccountID, amount)
	case domain.TransactionTypeIncome:
		negativeAmount := amount.Neg()
		return adapters.AccountRepo.UpdateBalance(*toAccountID, negativeAmount)
	case domain.TransactionTypeTransfer:
		err := adapters.AccountRepo.UpdateBalance(*fromAccountID, amount)
		if err != nil {
			return err
		}

		negativeAmount := amount.Neg()
		return adapters.AccountRepo.UpdateBalance(*toAccountID, negativeAmount)
	default:
		return domain.ErrInvalidTransactionType
	}
}
