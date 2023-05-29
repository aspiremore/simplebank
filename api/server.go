package api

import (
	db "github/aspiremore/simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// All HTTP server request goes through this struct
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)
	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(e error) gin.H {
	return gin.H{"error": e.Error()}
}
