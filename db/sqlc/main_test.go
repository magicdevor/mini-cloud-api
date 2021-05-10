package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"xiebeitech.com/mini-cloud-api/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
		config, err := util.NewConfig("../../.")
	if err != nil {
		log.Fatalln("cannot init config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBDataSource)
	if err != nil {
		log.Fatal("cannot connect database:", err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
