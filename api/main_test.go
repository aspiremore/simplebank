package api

import (
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func NewTestServer(t *testing.T, store db.Store) *Server  {
	config := util.Config{
		AccessTokenDuration: time.Minute,
		TokenSymmetricKey: util.RandomString(32),
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)
	return server
}
func TestMain(m *testing.M)  {
	gin.SetMode(gin.DebugMode)
	os.Exit(m.Run())
}
