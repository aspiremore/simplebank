package tests

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github/aspiremore/simplebank/db/sqlc"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)
var testDB *sql.DB
var testQueries *db.Queries

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("We cannot connect to the DB:", err)
	}
	testQueries = db.New(testDB)
	os.Exit(m.Run())

}

