package uow

import (
	"context"
	"database/sql"
)

const TxContextKey = "tx"

type Transaction interface {
	Begin(context.Context) (context.Context, error)
	Commit() error
	Rollback() error
}

type Base struct {
	db *sql.DB
}

func NewBase(db *sql.DB) Base {
	return Base{
		db: db,
	}
}

type uow struct {
	db *sql.DB
	tx *sql.Tx
}

func NewInstance(base Base) Transaction {
	return &uow{
		db: base.db,
	}
}

func (u *uow) Begin(ctx context.Context) (context.Context, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	u.tx = tx
	ctx = context.WithValue(ctx, TxContextKey, tx)

	return ctx, nil
}

func (u *uow) Commit() error {
	return u.tx.Commit()
}

func (u *uow) Rollback() error {
	return u.tx.Rollback()
}
