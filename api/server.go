package api

import (
	"fmt"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"github.com/aspiremore/simplebank/token"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// All HTTP server request goes through this struct
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config util.Config
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create paseto token maker")
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.createUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(e error) gin.H {
	return gin.H{"error": e.Error()}
}
