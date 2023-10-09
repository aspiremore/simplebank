package tests

import (
	"database/sql"
	"github/aspiremore/simplebank/db/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB
var testQueries *db.Queries

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("unable to load configurations", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("We cannot connect to the DB:", err)
	}
	testQueries = db.New(testDB)
	os.Exit(m.Run())

}

