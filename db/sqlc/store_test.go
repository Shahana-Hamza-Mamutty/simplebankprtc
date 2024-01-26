package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	testStore := NewStore(testDB)
	account1, _, _ := createTestAccount()
	account2, _, _ := createTestAccount()

	n := 1
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	fmt.Println("===========acc bal bef==========", account1.Balance, account2.Balance)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx: %d", i+1)
		ctx := context.WithValue(context.Background(), txKey, txName)
		go func() {
			res, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- res

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		res := <-results
		require.NoError(t, err)
		require.NotEmpty(t, res)
		// check transfer
		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		// fmt.Println("aaa==================")

		// check accounts
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		fmt.Println("===========acc bal tx==========", fromAccount.Balance, toAccount.Balance)

		// check balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 2)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
	}

	// check the final balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("===========acc bal fin==========", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, updatedAccount1.Balance, account1.Balance-(int64(n)*amount))
	require.Equal(t, updatedAccount2.Balance, account2.Balance+(int64(n)*amount))

}
