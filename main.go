package main

import (
	"database/sql"
	"log"

	"xiebeitech.com/mini-cloud-api/api"
	db "xiebeitech.com/mini-cloud-api/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/mini_db?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln("cannot connect database:", err)
	}
	queries := db.New(conn)
	server := api.NewServer(queries)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}