package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"xiebeitech.com/mini-cloud-api/api"
	db "xiebeitech.com/mini-cloud-api/db/sqlc"
	"xiebeitech.com/mini-cloud-api/util"
)

func main() {
	config, err := util.NewConfig(".")
	if err != nil {
		log.Fatalln("cannot init config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBDataSource)
	if err != nil {
		log.Fatalln("cannot connect database:", err)
	}

	queries := db.New(conn)
	server, err := api.NewServer(queries, config)

	if err != nil {
		log.Fatalln("cannot create a server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}
