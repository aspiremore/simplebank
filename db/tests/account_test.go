package tests

import (
	"context"
	"database/sql"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateRandomAccount(t *testing.T)  db.Account {
	user := CreateRandomUser(t)
	arg := db.CreateAccountParams{Balance: util.RandomMoney(), Owner: user.Username, Currency: util.RandomCurrency()}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Owner, account.Owner)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount( t *testing.T)  {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T){
	account := CreateRandomAccount(t)
	same_account, err := testQueries.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, same_account)
	require.Equal(t, account.Balance, same_account.Balance)
	require.Equal(t, account.Currency, same_account.Currency)
	require.Equal(t, account.Owner, same_account.Owner)
}

func TestUpdateAccount(t *testing.T){
	account := CreateRandomAccount(t)
	same_account, err := testQueries.UpdateAccount(context.Background(), db.UpdateAccountParams{ID : account.ID, Balance: 45})
	require.NoError(t, err)
	require.Equal(t, account.ID, same_account.ID)
	require.Equal(t, account.Currency, same_account.Currency)
	require.Equal(t, same_account.Balance, int64(45))
}


func TestDeleteAccount(t *testing.T){
	account := CreateRandomAccount(t)
	ctx := context.Background()
	err := testQueries.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)

	same_account, err := testQueries.GetAccount(ctx, account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, same_account)
}

