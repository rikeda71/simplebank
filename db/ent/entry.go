package db

import (
	"context"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/entry"
)

type CreateEntryParams struct {
	AccountID int `json:"account_id"`
	Amount    int `json:"amount"`
}

func (store *SQLStore) CreateEntry(ctx context.Context, tx *ent.Tx, arg CreateEntryParams) (*ent.Entry, error) {
	return tx.Entry.Create().
		SetAccountID(arg.AccountID).
		SetAmount(arg.Amount).
		Save(ctx)
}

func (store *SQLStore) GetEntry(ctx context.Context, id int) (*ent.Entry, error) {
	return store.entClient.Entry.Get(ctx, id)
}

type ListEntriesParams struct {
	AccountID int `json:"account_id"`
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
}

func (store *SQLStore) ListEntries(ctx context.Context, arg ListEntriesParams) ([]*ent.Entry, error) {
	return store.entClient.Entry.Query().
		Where(entry.AccountID(arg.AccountID)).
		Limit(arg.Limit).
		Offset(arg.Offset).
		All(ctx)
}
