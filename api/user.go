package api

import (
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var r CreateUserRequest
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	hashPassword, err := util.HashPassword(r.Password)
	if err !=  nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username:    r.Username,
		HashedPassword: hashPassword,
		Email:  r.Email,
		FullName: r.FullName,
	})
	if err != nil {
		if pqErr,ok:= err.(*pq.Error); ok {
			switch pqErr.Code {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	userResponse := db.UserResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
	ctx.JSON(http.StatusOK, userResponse)
}


