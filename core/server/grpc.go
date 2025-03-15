package server

import (
	"fmt"
	"github.com/alpha-omega-corp/cloud/core/database"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewGRPC(host string, dbHandler *database.Handler, proto func(db *bun.DB, grpc *grpc.Server)) error {
	listen, err := net.Listen("tcp", host)

	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	if dbHandler != nil {
		db := dbHandler.Database()
		defer func(db *bun.DB) {
			err := db.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(db)

		proto(db, srv)
	} else {
		proto(nil, srv)
	}

	fmt.Printf("running at tcp://%v", host)
	return srv.Serve(listen)
}
