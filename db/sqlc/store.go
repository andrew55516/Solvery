package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	PaymentTx(ctx context.Context, arg PaymentTxParams) (PaymentTxResult, error)
}

type PSQLStore struct {
	*Queries
	db *sql.DB
}

type PaymentTxParams struct {
	UserEmail string `json:"user_email"`
	Amount    int32  `json:"amount"`
	Comment   string `json:"comment"`
}

type PaymentTxResult struct {
	User  User  `json:"user"`
	Entry Entry `json:"entry"`
}

func NewStore(db *sql.DB) Store {
	return &PSQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *PSQLStore) PaymentTx(ctx context.Context, arg PaymentTxParams) (PaymentTxResult, error) {
	var res PaymentTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		res.User, err = q.AddUserCredit(ctx, AddUserCreditParams{
			Amount: arg.Amount,
			Email:  arg.UserEmail,
		})
		if err != nil {
			return err
		}

		res.Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			UserEmail: arg.UserEmail,
			Amount:    arg.Amount,
			Comment:   arg.Comment,
		})
		if err != nil {
			return err
		}

		return err
	})

	return res, err
}

func (store *PSQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
