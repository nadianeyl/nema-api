package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/govalues/decimal"

	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
	"github.com/nadianeyl/nema-api/internal/validator"
)

type TransactionService struct {
	TxProvider      *repository.TxProvider
	TransactionRepo repository.TransactionRepository
}

func NewTransactionService(txProvider *repository.TxProvider, transactionRepo repository.TransactionRepository) TransactionService {
	return TransactionService{
		TxProvider:      txProvider,
		TransactionRepo: transactionRepo,
	}
}

func (s *TransactionService) List(ctx context.Context, req *ListTransactionsRequest) ([]*TransactionResponse, domain.Metadata, error) {
	transactions, metadata, err := s.TransactionRepo.GetAllForUser(ctx, req.UserID, req.TransactionFilters)
	if err != nil {
		return nil, domain.Metadata{}, err
	}

	res := make([]*TransactionResponse, 0)
	for _, transaction := range transactions {
		res = append(res, &TransactionResponse{
			ID:            transaction.ID,
			UserID:        transaction.UserID,
			Type:          transaction.Type,
			CategoryID:    transaction.CategoryID,
			Amount:        transaction.Amount,
			Date:          transaction.Date,
			Title:         transaction.GetTitle(),
			Notes:         transaction.GetNotes(),
			FromAccountID: transaction.GetFromAccountID(),
			ToAccountID:   transaction.GetToAccountID(),
			CreatedAt:     transaction.CreatedAt,
			UpdatedAt:     transaction.UpdatedAt,
		})
	}

	return res, metadata, nil
}

func (s *TransactionService) Add(ctx context.Context, req *AddTransactionRequest) (*TransactionResponse, error) {
	var result *domain.Transaction

	err := s.TxProvider.WithTx(ctx, func(adapters repository.Adapters) error {
		_, err := adapters.CategoryRepo.GetByIDAndTypeForUser(ctx, req.CategoryID, req.Type, req.UserID)
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
			fromAccount, err = adapters.AccountRepo.GetByIDForUpdate(ctx, *req.FromAccountID)
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
			toAccount, err = adapters.AccountRepo.GetByIDForUpdate(ctx, *req.ToAccountID)
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

		err = adapters.TransactionRepo.Insert(ctx, transaction)
		if err != nil {
			return err
		}

		err = s.updateAccountBalances(ctx, adapters, transaction)
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

func (s *TransactionService) GetDetailByID(ctx context.Context, req *GetTransactionDetailRequest) (*TransactionDetailResponse, error) {
	detail, err := s.TransactionRepo.GetDetailByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if detail.Transaction.UserID != req.UserID {
		return nil, domain.ErrUserNotAllowed
	}

	res := &TransactionDetailResponse{
		Transaction: TransactionResponse{
			ID:            detail.Transaction.ID,
			UserID:        detail.Transaction.UserID,
			Type:          detail.Transaction.Type,
			CategoryID:    detail.Transaction.CategoryID,
			Amount:        detail.Transaction.Amount,
			Date:          detail.Transaction.Date,
			Title:         detail.Transaction.GetTitle(),
			Notes:         detail.Transaction.GetNotes(),
			FromAccountID: detail.Transaction.GetFromAccountID(),
			ToAccountID:   detail.Transaction.GetToAccountID(),
			CreatedAt:     detail.Transaction.CreatedAt,
			UpdatedAt:     detail.Transaction.UpdatedAt,
		},
		Category:    CategoryInfo(detail.Category),
		FromAccount: (*AccountInfo)(detail.FromAccount),
		ToAccount:   (*AccountInfo)(detail.ToAccount),
	}

	return res, nil
}

func (s *TransactionService) Update(ctx context.Context, req *UpdateTransactionRequest, v *validator.Validator) (*TransactionResponse, error) {
	var result *domain.Transaction

	err := s.TxProvider.WithTx(ctx, func(adapters repository.Adapters) error {
		transaction, err := adapters.TransactionRepo.GetByID(ctx, req.ID)
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

		_, err = adapters.CategoryRepo.GetByIDAndTypeForUser(ctx, transaction.CategoryID, transaction.Type, req.UserID)
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
			account, err := adapters.AccountRepo.GetByIDForUpdate(ctx, accountID)
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

		err = s.revertAccountBalances(ctx, adapters, oldType, oldAmount, oldFromAccountID, oldToAccountID)
		if err != nil {
			return err
		}

		switch transaction.Type {
		case domain.TransactionTypeExpense:
			transaction.SetToAccountID(nil)
		case domain.TransactionTypeIncome:
			transaction.SetFromAccountID(nil)
		}

		err = adapters.TransactionRepo.Update(ctx, transaction)
		if err != nil {
			return err
		}

		err = s.updateAccountBalances(ctx, adapters, transaction)
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

func (s *TransactionService) Delete(ctx context.Context, req *DeleteTransactionRequest) error {
	err := s.TxProvider.WithTx(ctx, func(adapters repository.Adapters) error {
		transaction, err := adapters.TransactionRepo.GetByIDForUpdate(ctx, req.ID)
		if err != nil {
			return err
		}

		if transaction.UserID != req.UserID {
			return domain.ErrUserNotAllowed
		}

		err = s.revertAccountBalances(ctx, adapters, transaction.Type, transaction.Amount, transaction.GetFromAccountID(), transaction.GetToAccountID())
		if err != nil {
			return err
		}

		err = adapters.TransactionRepo.Delete(ctx, transaction.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *TransactionService) updateAccountBalances(ctx context.Context, adapters repository.Adapters, transaction *domain.Transaction) error {
	switch transaction.Type {
	case domain.TransactionTypeExpense:
		negativeAmount := transaction.Amount.Neg()
		fromAccountID := transaction.GetFromAccountID()
		return adapters.AccountRepo.UpdateBalance(ctx, *fromAccountID, negativeAmount)
	case domain.TransactionTypeIncome:
		toAccountID := transaction.GetToAccountID()
		return adapters.AccountRepo.UpdateBalance(ctx, *toAccountID, transaction.Amount)
	case domain.TransactionTypeTransfer:
		negativeAmount := transaction.Amount.Neg()
		fromAccountID := transaction.GetFromAccountID()
		toAccountID := transaction.GetToAccountID()

		err := adapters.AccountRepo.UpdateBalance(ctx, *fromAccountID, negativeAmount)
		if err != nil {
			return err
		}

		return adapters.AccountRepo.UpdateBalance(ctx, *toAccountID, transaction.Amount)
	default:
		return domain.ErrInvalidTransactionType
	}
}

func (s *TransactionService) revertAccountBalances(ctx context.Context, adapters repository.Adapters, transactionType domain.TransactionType, amount decimal.Decimal, fromAccountID, toAccountID *uuid.UUID) error {
	switch transactionType {
	case domain.TransactionTypeExpense:
		return adapters.AccountRepo.UpdateBalance(ctx, *fromAccountID, amount)
	case domain.TransactionTypeIncome:
		negativeAmount := amount.Neg()
		return adapters.AccountRepo.UpdateBalance(ctx, *toAccountID, negativeAmount)
	case domain.TransactionTypeTransfer:
		err := adapters.AccountRepo.UpdateBalance(ctx, *fromAccountID, amount)
		if err != nil {
			return err
		}

		negativeAmount := amount.Neg()
		return adapters.AccountRepo.UpdateBalance(ctx, *toAccountID, negativeAmount)
	default:
		return domain.ErrInvalidTransactionType
	}
}
