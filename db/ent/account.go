package db

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/account"
)

type AddAccountBalanceParams struct {
	Amount int `json:"amount"`
	ID     int `json:"id"`
}

// AddAccountBalance update account balance
func (store *SQLStore) AddAccountBalance(ctx context.Context, tx *ent.Tx, arg AddAccountBalanceParams) (*ent.Account, error) {
	account, err := store.GetAccountForUpdate(ctx, tx, arg.ID)
	if err != nil {
		return nil, err
	}
	return account.Update().AddBalance(arg.Amount).Save(ctx)
}

// CRUD

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}

func (store *SQLStore) CreateAccount(ctx context.Context, arg CreateAccountParams) (*ent.Account, error) {
	return store.entClient.Account.Create().
		SetOwner(arg.Owner).
		SetBalance(arg.Balance).
		SetCurrency(arg.Currency).
		Save(context.Background())
}

func (store *SQLStore) DeleteAccount(ctx context.Context, id int) error {
	return store.entClient.Account.DeleteOneID(id).Exec(ctx)
}

func (store *SQLStore) GetAccount(ctx context.Context, id int) (*ent.Account, error) {
	return store.entClient.Account.Get(ctx, id)
}

func (store *SQLStore) GetAccountForUpdate(ctx context.Context, tx *ent.Tx, id int) (*ent.Account, error) {
	// SELECT FOR UPDATE を使って一貫性を保つ
	// Get だと通常の SELECT になる
	// account1, err := tx.Account.Get(ctx, arg.FromAccountID)
	// if err != nil {
	// 	return err
	// }
	// result.ToAccount, err = tx.Account.UpdateOneID(account1.ID).
	// 	AddBalance(-arg.Ammount).
	// 	Save(ctx)
	return tx.Account.Query().
		// `FOR UPDATE` を使う場合はシンプルに以下のように記述できる
		// Where(account.ID(arg.FromAccountID)).
		// ForUpdate().

		// FOR NO KEY UPDATE はビルダーメソッドが用意されていないので下記のようにして書く
		Where(func(s *sql.Selector) {
			s.Where(sql.EQ(account.FieldID, id)).
				For(sql.LockNoKeyUpdate)
		}).
		// Unique(false) を呼び出すことで SELECT DISTINCT => SELECT に変更
		// psql において、FOR NO KEY UPDATE は SELECT DISTINCT と両立できない
		Unique(false).
		Only(ctx)
}

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func (store *SQLStore) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]*ent.Account, error) {
	return store.entClient.Account.Query().
		Where(account.OwnerContains(arg.Owner)).
		Limit(arg.Limit).
		Offset(arg.Offset).
		All(ctx)
}

type UpdateAccountParams struct {
	ID      int `json:"id"`
	Balance int `json:"balance"`
}

func (store *SQLStore) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (*ent.Account, error) {
	return store.entClient.Account.UpdateOneID(arg.ID).
		SetBalance(arg.Balance).
		Save(ctx)
}
