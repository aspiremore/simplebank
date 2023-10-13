package tests

import (
	"context"
	"fmt"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := db.NewStore(testDB)
	from_account := CreateRandomAccount(t)
	to_account := CreateRandomAccount(t)
	errs := make(chan error)
	results := make(chan db.TransferTxResults)
	amount := int64(10)
	n := 16

	fmt.Println("before >>>: ", from_account.Balance, to_account.Balance)
	for i := 0; i<n; i++ {

		go func() {
			txName := fmt.Sprintf("transaction %d", i+1)
			ctx := context.WithValue(context.Background(), db.TxKey , txName)
			txResult, err := store.TransferTx(ctx, db.TransferTxParams{
				FromAccountID: from_account.ID,
				ToAccountID: to_account.ID,
				Amount: amount,
			})
			errs <- err
			results <- txResult
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <- errs
		result := <- results
		require.Empty(t, err)

		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, from_account.ID)
		require.Equal(t, transfer.ToAccountID, to_account.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)


		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.Empty(t, err)


		//check entries
		from_entry := result.FromEntry
		require.NotZero(t, from_entry)
		require.NotZero(t, from_entry.ID)
		require.NotZero(t, from_entry.CreatedAt)
		require.Equal(t, from_entry.Amount, -amount)
		require.Equal(t, from_entry.AccountID, from_account.ID)

		_, err = store.GetEntry(context.Background(), from_entry.ID)
		require.NoError(t, err)


		//check entries
		to_entry := result.ToEntry
		require.NotZero(t, to_entry)
		require.NotZero(t, to_entry.ID)
		require.NotZero(t, to_entry.CreatedAt)
		require.Equal(t, to_entry.Amount, amount)
		require.Equal(t, to_entry.AccountID, to_account.ID)

		_, err = store.GetEntry(context.Background(), to_entry.ID)
		require.NoError(t, err)

		//todo account balance update
		fromAccount := result.FromAccount
		require.NotZero(t, fromAccount.ID)
		require.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		require.NotZero(t, toAccount.ID)
		require.NotEmpty(t, toAccount)

		//check account balance
		fmt.Println("txn >>> :", fromAccount.Balance, toAccount.Balance, amount)
		diff1 := from_account.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - to_account.Balance
		require.Equal(t, diff2, diff1)
		require.True(t, diff1 % amount == 0)

	}
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), from_account.ID)
	require.Empty(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), to_account.ID)
	require.Empty(t, err)

	fmt.Println("After >>> : ", updatedFromAccount.Balance, updatedToAccount.Balance)
}


func TestTransferTxDeadlock(t *testing.T) {
	store := db.NewStore(testDB)
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	errs := make(chan error)

	amount := int64(10)
	n := 10

	fmt.Println("before >>>: ", account1.Balance, account2.Balance)
	for i := 0; i<n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			txName := fmt.Sprintf("transaction %d", i+1)
			ctx := context.WithValue(context.Background(), db.TxKey , txName)
			_, err := store.TransferTx(ctx, db.TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID: toAccountID,
				Amount: amount,
			})
			errs <- err
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <- errs
		require.Empty(t, err)

	}
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Empty(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.Empty(t, err)

	fmt.Println("After >>> : ", updatedFromAccount.Balance, updatedToAccount.Balance)
	require.Equal(t, updatedFromAccount.Balance, account1.Balance)
	require.Equal(t, updatedToAccount.Balance, account2.Balance )
}
