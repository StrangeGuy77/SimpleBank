package db

import (
	"context"
	"testing"

	"github.com/StrangeGuy77/SimpleBank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := make(chan TransferTxParams)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run concurrent transactions
	for range 5 {
		go func() {
			arg := TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        util.RandomMoney(),
			}

			result, err := store.TransferTx(context.Background(), arg)

			args <- arg
			errs <- err
			results <- result
		}()
	}

	for range 5 {
		arg := <-args
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
	}
}
