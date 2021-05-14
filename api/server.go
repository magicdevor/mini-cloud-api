package api

import (
	"github.com/gin-gonic/gin"
	db "xiebeitech.com/mini-cloud-api/db/sqlc"
	"xiebeitech.com/mini-cloud-api/token"
	"xiebeitech.com/mini-cloud-api/util"
)

type Server struct {
	config    util.Config
	db        *db.Queries
	router    *gin.Engine
	snowflake *util.Snowflake
	maker     token.Maker
}

func NewServer(db *db.Queries, cf util.Config) (server *Server, err error) {
	server = &Server{
		config:    cf,
		db:        db,
		snowflake: util.NewSnowflake(123),
	}

	tokenMaker, err := token.NewJWTMaker(cf.TokenSymmetricKey)
	if err != nil {
		return
	}

	server.maker = tokenMaker

	router := gin.Default()
	v2 := router.Group("/v2", JWTAuthenticateMiddleware(tokenMaker))
	{
		v2.GET("/profile", server.GetProfile)
		v2.POST("/profile", server.CreateProfile)
	}
	router.GET("/login", AuthenticateMiddleware(server.config), server.Login)
	server.router = router

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
