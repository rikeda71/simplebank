package db

import (
	"context"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/transfer"
)

type CreateTransferParams struct {
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	Amount        int `json:"amount"`
}

func (store *SQLStore) CreateTransfer(ctx context.Context, tx *ent.Tx, arg CreateTransferParams) (*ent.Transfer, error) {
	return tx.Transfer.Create().
		SetFromAccountID(arg.FromAccountID).
		SetToAccountID(arg.ToAccountID).
		SetAmount(arg.Amount).
		Save(ctx)
}

func (store *SQLStore) GetTransfer(ctx context.Context, id int) (*ent.Transfer, error) {
	return store.entClient.Transfer.Get(ctx, id)
}

type ListTransfersParams struct {
	FromAccountID int `json:"from_account_id"`
	ToAccountID   int `json:"to_account_id"`
	Limit         int `json:"limit"`
	Offset        int `json:"offset"`
}

func (store *SQLStore) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]*ent.Transfer, error) {
	return store.entClient.Transfer.Query().
		Where(transfer.FromAccountID(arg.FromAccountID), transfer.ToAccountID(arg.ToAccountID)).
		Limit(arg.Limit).
		Offset(arg.Offset).
		All(ctx)
}
