package api

import (
	"fmt"
	mockdb "github.com/aspiremore/simplebank/db/mock"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func randomAccount(t *testing.T) db.Account {
	return db.Account{
		ID: util.RandomInt(1,100),
		Owner: util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance: util.RandomMoney(),
	}
}

func TestGetAccount(t *testing.T)  {

	account := randomAccount(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetAccount(gomock.Any(),gomock.Eq(account.ID)).
		Times(1).
		Return(account,nil)
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d",account.ID)
	request, err := http.NewRequest(http.MethodGet,url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder,request)
	require.Equal(t,http.StatusOK, recorder.Code)
}
