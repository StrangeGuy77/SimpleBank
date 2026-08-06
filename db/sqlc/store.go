package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Instance a Transaction Operation Function Wrapper
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Instance Queries using the Transaction Wrapper
	q := New(tx)
	// Call the callback function using the queries with the Wrapper
	err = fn(q)

	if err != nil {
		// If callback returns an error, rollback the whole transaction.
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("Unable to rollback transaction. TX err %v, Rollback err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Creates a transfer record, add account entries, and update accounts balances within a single DB transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create transfer record
		result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		// Create entry telling the money is being withdrawn from the FromAccount Account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// Create entry telling the money is being moving into the ToAccount Account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		// TODO: update accounts balance

		return nil
	})

	return result, err
}
