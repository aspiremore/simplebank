package api

import (
	"database/sql"
	"fmt"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
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


type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	access_token string `json:"access___token"`
	user db.UserResponse `json:"user"`
}

func  NewUserResponse(user db.User) db.UserResponse {
	return db.UserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}

}


func (s *Server) loginUser(ctx *gin.Context)  {
	fmt.Println("----------------------------------")
	var loginUserRequest LoginUserRequest
	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	u, err := s.store.GetUser(ctx,loginUserRequest.Username)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return

	}
	err = util.CheckPassword(u.HashedPassword,loginUserRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	//duration will come from viper in a moment
	token,err := s.tokenMaker.CreateToken(loginUserRequest.Username,time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	response := LoginUserResponse{
		access_token: token,
		user: NewUserResponse(u),
	}
	fmt.Printf("--------------------dfsdfsdfsdf--------------%+v", response)
	ctx.JSON(http.StatusOK,response)
}



