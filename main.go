package main

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/nhatnhanchiha/simpleBank/api"
	db "github.com/nhatnhanchiha/simpleBank/db/sqlc"
	"github.com/nhatnhanchiha/simpleBank/gapi"
	"github.com/nhatnhanchiha/simpleBank/pb"
	"github.com/nhatnhanchiha/simpleBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
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
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	c, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	err = pb.RegisterSimpleBankHandlerServer(c, grpcMux, server)
	if err != nil {
		log.Fatal("Cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server")
	}

}

func runGinServer(err error, config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot Create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
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
