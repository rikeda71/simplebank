package db

import (
	"context"
	"fmt"

	"github.com/s14t284/simplebank/ent"
)

type Store struct {
	client *ent.Client
}

func NewStore(client *ent.Client) *Store {
	return &Store{
		client: client,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, err := store.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("tx err: %v, rb err: %v", err, rerr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	Ammount       int `json:"ammount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    *ent.Transfer `json:"transfer"`
	FromAccount *ent.Account  `json:"from_account_id"`
	ToAccount   *ent.Account  `json:"to_account_id"`
	FromEntry   *ent.Entry    `json:"from_entry"`
	ToEntry     *ent.Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, and account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(tx *ent.Tx) error {
		var err error
		result.Transfer, err = tx.Transfer.Create().
			SetFromAccountID(arg.FromAccountID).
			SetToAccountID(arg.ToAccountID).
			SetAmount(arg.Ammount).
			Save(ctx)
		if err != nil {
			return err
		}

		// From
		result.FromEntry, err = tx.Entry.Create().
			SetAccountID(arg.FromAccountID).
			SetAmount(-arg.Ammount).
			Save(ctx)
		if err != nil {
			return err
		}

		// To
		result.ToEntry, err = tx.Entry.Create().
			SetAccountID(arg.ToAccountID).
			SetAmount(arg.Ammount).
			Save(ctx)
		if err != nil {
			return err
		}

		// TODO: update accounts' balance

		return nil
	})

	return result, err
}
