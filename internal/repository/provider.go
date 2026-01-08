package repository

import (
	"database/sql"

	"github.com/nadianeyl/nema-api/internal/helper"
)

type txFunc func(adapters Adapters) error

type Adapters struct {
	CategoryRepo    CategoryRepository
	AccountRepo     AccountRepository
	TransactionRepo TransactionRepository
}

type TxProvider struct {
	DB *sql.DB
}

func NewTxProvider(db *sql.DB) *TxProvider {
	return &TxProvider{DB: db}
}

func (p *TxProvider) WithTx(txFn txFunc) error {
	return helper.RunInTx(p.DB, func(tx *sql.Tx) error {
		adapters := Adapters{
			CategoryRepo:    CategoryRepository{DB: tx},
			AccountRepo:     AccountRepository{DB: tx},
			TransactionRepo: TransactionRepository{DB: tx},
		}

		return txFn(adapters)
	})
}
