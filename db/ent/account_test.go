package db

import (
	"context"
	"testing"
	"time"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) *ent.Account {
	client := createEntClient(t)
	defer client.Close()

	owner := util.RandomString(3)
	balance := int(util.RandomMoney())
	currency := util.RandomCurrency()
	account, err := client.Account.Create().
		SetOwner(owner).
		SetBalance(balance).
		SetCurrency(currency).
		Save(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, owner, account.Owner)
	require.Equal(t, balance, account.Balance)
	require.Equal(t, currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	account1 := createRandomAccount(t)
	account2, err := client.Account.Get(context.Background(), account1.ID)
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

	account1 := createRandomAccount(t)

	balance := int(util.RandomMoney())
	account2, err := client.Account.UpdateOneID(account1.ID).SetBalance(balance).Save(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	account1 := createRandomAccount(t)
	err := client.Account.DeleteOneID(account1.ID).Exec(context.Background())
	require.NoError(t, err)

	account2, err := client.Account.Get(context.Background(), account1.ID)
	require.Error(t, err)
	require.Equal(t, err.Error(), "ent: account not found")
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	client := createEntClient(t)
	defer client.Close()

	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := client.Account.Query().All(context.Background())
	require.NoError(t, err)
	require.Len(t, accounts, 10)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
