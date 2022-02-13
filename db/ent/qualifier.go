package db

import (
	"context"

	"github.com/s14t284/simplebank/ent"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (*ent.Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (*ent.Account, error)
	// CreateEntry(ctx context.Context, arg CreateEntryParams) (*ent.Entry, error)
	// CreateTransfer(ctx context.Context, arg CreateTransferParams) (*ent.Transfer, error)
	DeleteAccount(ctx context.Context, id int) error
	GetAccount(ctx context.Context, id int) (*ent.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (*ent.Account, error)
	GetEntry(ctx context.Context, id int) (*ent.Entry, error)
	GetTransfer(ctx context.Context, id int) (*ent.Transfer, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]*ent.Account, error)
	// ListEntries(ctx context.Context, arg ListEntriesParams) ([]*ent.Entry, error)
	// ListTransfers(ctx context.Context, arg ListTransfersParams) ([]*ent.Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (*ent.Account, error)
}
