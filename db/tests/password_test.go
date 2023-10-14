package tests

import (
	"github.com/aspiremore/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password := util.RandomString(7)
	hashedpassword,err := util.HashPassword(password)
	require.NoError(t ,err)
	require.NotEmpty(t, password)

	err = util.CheckPassword(hashedpassword,password)
	require.NoError(t ,err)

	wrongPassword := util.RandomString(7)
	err = util.CheckPassword(wrongPassword,password)
	require.Error(t ,err, bcrypt.ErrMismatchedHashAndPassword.Error())
}


