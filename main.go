package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// run db migration
	runDbMigration(config.MigrationURL, config.DbSource)

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

func runDbMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	if err = migration.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("err no change")
		} else {
			log.Fatal(err)
		}
	}

	log.Println("db migrated successfully")
}
