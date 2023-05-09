package api

import (
	"fmt"
	db "github/aspiremore/simplebank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
    Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) CreateAccount(ctx *gin.Context) {
	var r CreateAccountRequest
	if err := ctx.ShouldBindJSON(r); err != nil {
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
