package token

import (
	"github.com/aspiremore/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)
//valid case
func TestJWTMaker(t *testing.T)  {
	tokenMaker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	//create token
	username := util.RandomOwner()
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//validate token
	payload, err := tokenMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t,expiredAt, payload.ExpiredAt, time.Second)
	require.WithinDuration(t,issuedAt, payload.IssuedAt, time.Second)
}

func TestExpiredToken(t *testing.T)  {
	tokenMaker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	//create token
	username := util.RandomOwner()
	token, err := tokenMaker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)


	//validate token
	payload, err := tokenMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
