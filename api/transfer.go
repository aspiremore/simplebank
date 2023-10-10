package api

import (
	"database/sql"
	"fmt"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransferAccountRequest struct {
	FromAccountID    int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID    int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	var r TransferAccountRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	fmt.Println("-------------------------")

	if !s.validateAccount(ctx,r.FromAccountID,r.Currency){
		return
	}
	if !s.validateAccount(ctx,r.ToAccountID,r.Currency){
		return
	}

	result, err := s.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID:    r.FromAccountID,
		ToAccountID: r.ToAccountID,
		Amount:  r.Amount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (s *Server)validateAccount(ctx *gin.Context, account_id int64, currency string) bool {
	account, err := s.store.GetAccount(ctx,account_id)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s",account_id,account.Currency, currency)
		ctx.JSON(http.StatusBadRequest,errResponse(err))
		return false
	}
	return true
}


