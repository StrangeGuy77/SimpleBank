package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := TransferTxParams{
		FromAccountId: account1.ID,
		ToAccountId:   account2.ID,
		Amount:        10,
	}

	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 5

	// run concurrent transactions
	for range n {
		go func() {
			result, err := store.TransferTx(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for range n {
		err := <-errs

		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Validate transfer
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, arg.FromAccountId, result.Transfer.FromAccountID)
		require.Equal(t, arg.ToAccountId, result.Transfer.ToAccountID)
		require.Equal(t, arg.Amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, tErr := testQueries.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, tErr)

		// Validate entries
		// FROM
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, arg.FromAccountId, result.FromEntry.AccountID)
		require.Equal(t, -arg.Amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.CreatedAt)
		require.NotZero(t, result.FromEntry.ID)

		_, ftErr := testQueries.GetEntry(context.Background(), result.FromEntry.ID)
		require.NoError(t, ftErr)

		// TO
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, arg.ToAccountId, result.ToEntry.AccountID)
		require.Equal(t, arg.Amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.CreatedAt)
		require.NotZero(t, result.ToEntry.ID)

		_, ttErr := testQueries.GetEntry(context.Background(), result.ToEntry.ID)
		require.NoError(t, ttErr)

		// Check accounts
		require.NotEmpty(t, result.FromAccount)
		require.Equal(t, result.FromAccount.ID, arg.FromAccountId)

		require.NotEmpty(t, result.ToAccount)
		require.Equal(t, result.ToAccount.ID, arg.ToAccountId)

		// Check balances
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%arg.Amount == 0) // 1 * amount, 2 * amount ...

		k := int(diff1 / arg.Amount)
		// K should always
		// be higher than 0 (otherwise something failed with transaction amount on account1 balance) and
		// lower than n (number of transactions being commited)
		require.True(t, k >= 0 && k <= n)
		require.NotContains(t, existed, k)

		existed[k] = true
	}

	// Check final balance of both accounts
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-(int64(n)*arg.Amount), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+(int64(n)*arg.Amount), updatedAccount2.Balance)
}
