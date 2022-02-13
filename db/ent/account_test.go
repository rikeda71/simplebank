package db

import (
	"context"
	"testing"
	"time"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, store *Store) *ent.Account {
	arg := CreateAccountParams{
		Owner:    util.RandomString(3),
		Balance:  int(util.RandomMoney()),
		Currency: util.RandomCurrency(),
	}

	account, err := store.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	store := NewStore(client)
	createRandomAccount(t, store)
}

func TestGetAccount(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	store := NewStore(client)

	account1 := createRandomAccount(t, store)
	account2, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	store := NewStore(client)

	account1 := createRandomAccount(t, store)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: int(util.RandomMoney()),
	}
	account2, err := store.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	store := NewStore(client)

	account1 := createRandomAccount(t, store)
	err := store.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := client.Account.Get(context.Background(), account1.ID)
	require.Error(t, err)
	require.Equal(t, err.Error(), "ent: account not found")
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	store := NewStore(client)

	for i := 0; i < 10; i++ {
		createRandomAccount(t, store)
	}

	accounts, err := store.ListAccounts(context.Background(), ListAccountsParams{Owner: "", Limit: 10, Offset: 0})
	require.NoError(t, err)
	require.Len(t, accounts, 10)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
