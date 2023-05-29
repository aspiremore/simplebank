package api

import (
	"database/sql"
	"fmt"
	db "github/aspiremore/simplebank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
    Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var r CreateAccountRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	fmt.Println("-------------------------")
	account, err := s.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    r.Owner,
		Currency: r.Currency,
		Balance:  0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}


type getAccountRequest struct {
	ID    int64 `uri:"id" binding:"required",min=1`
}

func (s *Server) getAccount(ctx *gin.Context) {
	var r getAccountRequest
	if err := ctx.ShouldBindUri(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	account, err := s.store.GetAccount(ctx, r.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
		return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}


type listAccountRequest struct {
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=1,max=5"`
}

func (s *Server) listAccount(ctx *gin.Context) {
	var r listAccountRequest
	if err := ctx.ShouldBindQuery(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := db.ListAccountsParams{
		Limit: r.PageSize,
		Offset: (r.PageID - 1) * r.PageSize,
	}
	fmt.Println("arg ", arg)
	accounts, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	fmt.Println("accounts ", accounts)
	ctx.JSON(http.StatusOK, accounts)
}
