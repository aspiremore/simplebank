package tests

import (
	"context"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateRandomUser(t *testing.T)  db.User {
	arg := db.CreateUserParams{
		Username: util.RandomOwner(),
		Email: util.RandomEmail(),
		HashedPassword: "secret",
		FullName: util.RandomOwner(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email,user.Email)

	require.NotZero(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser( t *testing.T)  {
	CreateRandomUser(t)
}

