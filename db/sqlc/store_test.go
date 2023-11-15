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

	exists := make(map[int]bool)
	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results
		require.NoError(t, err)
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

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)

		// check account's balance
		dif1 := account1.Balance - fromAccount.Balance
		dif2 := toAccount.Balance - account2.Balance
		// fmt.Printf("dif: %v %v", dif1, dif2)
		require.Equal(t, dif1, dif2)
		require.True(t, dif1 > 0)
		require.True(t, dif1%amount == 0)

		k := int(dif1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exists, k)
		exists[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updatedAccount1.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)

}
func TestTranferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer transactions
	n := 6
	amount := int64(10)

	errs := make(chan error)
	// results := make(chan TranferTxResult)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TranferTx(context.Background(), TranferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
			// results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		// result := <-results
		require.NoError(t, err)
	}
	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)

}

// func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {}
