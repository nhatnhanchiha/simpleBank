package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/nhatnhanchiha/simpleBank/api"
	db "github.com/nhatnhanchiha/simpleBank/db/sqlc"
	"github.com/nhatnhanchiha/simpleBank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot Create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
}
