package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/govalues/decimal"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

func (t TransactionType) String() string {
	return string(t)
}

func GetTransactionTypes() []TransactionType {
	return []TransactionType{TransactionTypeIncome, TransactionTypeExpense, TransactionTypeTransfer}
}

type Transaction struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Type          TransactionType
	CategoryID    uuid.UUID
	Amount        decimal.Decimal
	Date          time.Time
	Title         sql.NullString
	Notes         sql.NullString
	FromAccountID uuid.NullUUID
	ToAccountID   uuid.NullUUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Version       int
}

func (t *Transaction) GetTitle() *string {
	if t.Title.Valid {
		return &t.Title.String
	}
	return nil
}

func (t *Transaction) SetTitle(title *string) {
	if title != nil {
		t.Title = sql.NullString{String: *title, Valid: true}
	} else {
		t.Title = sql.NullString{}
	}
}

func (t *Transaction) GetNotes() *string {
	if t.Notes.Valid {
		return &t.Notes.String
	}
	return nil
}

func (t *Transaction) SetNotes(notes *string) {
	if notes != nil {
		t.Notes = sql.NullString{String: *notes, Valid: true}
	} else {
		t.Notes = sql.NullString{}
	}
}

func (t *Transaction) GetFromAccountID() *uuid.UUID {
	if t.FromAccountID.Valid {
		return &t.FromAccountID.UUID
	}
	return nil
}

func (t *Transaction) SetFromAccountID(id *uuid.UUID) {
	if id != nil {
		t.FromAccountID = uuid.NullUUID{UUID: *id, Valid: true}
	} else {
		t.FromAccountID = uuid.NullUUID{}
	}
}

func (t *Transaction) GetToAccountID() *uuid.UUID {
	if t.ToAccountID.Valid {
		return &t.ToAccountID.UUID
	}
	return nil
}

func (t *Transaction) SetToAccountID(id *uuid.UUID) {
	if id != nil {
		t.ToAccountID = uuid.NullUUID{UUID: *id, Valid: true}
	} else {
		t.ToAccountID = uuid.NullUUID{}
	}
}
