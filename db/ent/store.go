package db

import (
	"context"
	"fmt"

	"github.com/s14t284/simplebank/ent"
)

type Store struct {
	entClient *ent.Client
}

func NewStore(client *ent.Client) *Store {
	return &Store{
		entClient: client,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, err := store.entClient.Tx(ctx)
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

		// 同時に2つ以上のトランザクションが2つとも同一accountIdを指定し、FromとToが入れ替わっていた場合のデッドロック対策
		// accountId が小さい順に処理することでデッドロックを避けている
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = store.AddMoney(ctx, tx, arg.FromAccountID, -arg.Ammount, arg.ToAccountID, arg.Ammount)
		} else {
			result.ToAccount, result.FromAccount, err = store.AddMoney(ctx, tx, arg.ToAccountID, arg.Ammount, arg.FromAccountID, -arg.Ammount)
		}

		return nil
	})

	return result, err
}

func (store *Store) AddMoney(ctx context.Context, tx *ent.Tx, accountID1, amount1, accountID2, amount2 int) (account1, account2 *ent.Account, err error) {
	account1, err = store.AddAccountBalance(ctx, tx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = store.AddAccountBalance(ctx, tx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}
