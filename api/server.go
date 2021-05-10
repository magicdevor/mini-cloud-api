package api

import (
	"github.com/gin-gonic/gin"
	db "xiebeitech.com/mini-cloud-api/db/sqlc"
)

type Server struct {
	db     *db.Queries
	router *gin.Engine
}

func NewServer(db *db.Queries) *Server {
	server := &Server{
		db: db,
	}

	router := gin.Default()
	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
