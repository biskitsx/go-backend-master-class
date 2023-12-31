package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/biskitsx/go-backend-master-class/api"
	db "github.com/biskitsx/go-backend-master-class/db/sqlc"
	"github.com/biskitsx/go-backend-master-class/util"
	_ "github.com/golang/mock/mockgen/model"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	conn, err := sql.Open(config.DBdriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		panic(err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		panic(err)
	}

}
