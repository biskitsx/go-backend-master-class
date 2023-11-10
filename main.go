package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/biskitsx/go-backend-master-class/api"
	db "github.com/biskitsx/go-backend-master-class/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:root@localhost:5434/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(*store)
	err = server.Start(serverAddress)
	if err != nil {
		panic(err)
	}

}
