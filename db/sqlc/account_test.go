package db

import (
	"context"
	"log"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestAccount() (Account, *CreateAccountParams, error) {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  1000,
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	return account, &args, err

}
func TestCreateAccount(t *testing.T) {
	account, args, err := createTestAccount()
	log.Println(account)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Owner, args.Owner)
	require.Equal(t, account.Balance, args.Balance)
	require.Equal(t, account.Currency, args.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	account, _, _ := createTestAccount()
	accountGet, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, accountGet.Owner, account.Owner)
	require.Equal(t, accountGet.Balance, account.Balance)
	require.Equal(t, accountGet.Currency, account.Currency)
}

func TestUpdateAccount(t *testing.T) {
	account, _, _ := createTestAccount()
	args := UpdateAccountParams{
		ID:      account.ID,
		Balance: 1999,
	}
	accountUpdated, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)
	require.Equal(t, accountUpdated.Balance, args.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account, _, _ := createTestAccount()
	testQueries.DeleteAccount(context.Background(), account.ID)
	_, err1 := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err1)
}
