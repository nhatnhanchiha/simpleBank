package gapi

import (
	"fmt"
	db "github.com/nhatnhanchiha/simpleBank/db/sqlc"
	"github.com/nhatnhanchiha/simpleBank/pb"
	"github.com/nhatnhanchiha/simpleBank/token"
	"github.com/nhatnhanchiha/simpleBank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker")
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker}

	return server, nil
}
