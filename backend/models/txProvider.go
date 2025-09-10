package models

import (
	"database/sql"

	"github.com/jakubdrobny/speedcubingslovakia/backend/repository"
)

type TransactionProvider struct {
	db *sql.DB
}

func NewTransactionProvider(db *sql.DB) *TransactionProvider {
	return &TransactionProvider{
		db: db,
	}
}

func (p *TransactionProvider) Transact(txFunc func(adapters Adapters) error) error {
	return runInTx(p.db, func(tx *sql.Tx) error {
		adapters := repository.Adapters{
			User: NewUser(tx),
		}

		return txFunc(adapters)
	})
}
