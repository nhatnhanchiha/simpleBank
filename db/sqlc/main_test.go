package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nhatnhanchiha/simpleBank/util"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, _ := util.LoadConfig(".")
	var err error
	testDb, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
