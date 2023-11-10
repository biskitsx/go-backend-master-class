package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TranferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TranferTx(context.Background(), TranferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)
		require.Equal(t, result.Tranfer.Amount, amount)
		require.Equal(t, result.FromAccount.ID, account1.ID)
		require.Equal(t, result.ToAccount.ID, account2.ID)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.FromEntry.Amount)

		_, err = store.GetTransfer(context.Background(), result.Tranfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.Amount, -amount)
		require.Equal(t, fromEntry.AccountID, account1.ID)

		_, err = store.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.Amount, amount)
		// require.Equal(t, toEntry.AccountID, account2.ID)

		_, err = store.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)

	}
}

// func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {}
