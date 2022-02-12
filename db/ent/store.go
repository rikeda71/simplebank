package db

import (
	"context"
	"fmt"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/account"
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

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, and account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(tx *ent.Tx) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")

		result.Transfer, err = tx.Transfer.Create().
			SetFromAccountID(arg.FromAccountID).
			SetToAccountID(arg.ToAccountID).
			SetAmount(arg.Ammount).
			Save(ctx)
		if err != nil {
			return err
		}

		// From
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = tx.Entry.Create().
			SetAccountID(arg.FromAccountID).
			SetAmount(-arg.Ammount).
			Save(ctx)
		if err != nil {
			return err
		}

		// To
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = tx.Entry.Create().
			SetAccountID(arg.ToAccountID).
			SetAmount(arg.Ammount).
			Save(ctx)
		if err != nil {
			return err
		}

		// SELECT FOR UPDATE を使って一貫性を保つ
		// Get だと通常の SELECT になる
		// account1, err := tx.Account.Get(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		// result.ToAccount, err = tx.Account.UpdateOneID(account1.ID).
		// 	AddBalance(-arg.Ammount).
		// 	Save(ctx)
		fmt.Println(txName, "get account 1")
		q1, err := tx.Account.Query().
			Where(account.ID(arg.FromAccountID)).
			ForUpdate().
			Only(ctx)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update account 1")
		_, err = q1.Update().AddBalance(-arg.Ammount).Save(ctx)
		if err != nil {
			return err
		}

		// SELECT FOR UPDATE を使って一貫性を保つ
		// account2, err := tx.Account.Get(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }
		// result.ToAccount, err = tx.Account.UpdateOneID(account2.ID).
		// 	AddBalance(arg.Ammount).
		// 	Save(ctx)
		fmt.Println(txName, "get account 2")
		q2, err := tx.Account.Query().
			Where(account.ID(arg.ToAccountID)).
			ForUpdate().
			Only(ctx)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update account 2")
		_, err = q2.Update().AddBalance(arg.Ammount).Save(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
