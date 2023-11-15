package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/biskitsx/go-backend-master-class/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	testDB, err = sql.Open(config.DBdriver, config.DBSource)
	if err != nil {
		panic(err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
