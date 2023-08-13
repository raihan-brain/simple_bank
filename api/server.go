package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/raihan-brain/simple-bank/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

var router = gin.Default()

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}

	getRoutes(server)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func getRoutes(server *Server) {

	addAccountRoutes(router, server)
}
