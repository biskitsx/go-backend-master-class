package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5434/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	// conn, err := pgx.Connect(context.Background(), dbSource)
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err)
	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
