package api

import (
	"github.com/gin-gonic/gin"
	db "xiebeitech.com/mini-cloud-api/db/sqlc"
	"xiebeitech.com/mini-cloud-api/util"
)

type Server struct {
	config util.Config
	db     *db.Queries
	router *gin.Engine
}

func NewServer(db *db.Queries, cf util.Config) *Server {
	server := &Server{
		config: cf,
		db:     db,
	}

	router := gin.Default()
	router.GET("/login", AuthenticateMiddleware(server.config), server.Login)
	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
