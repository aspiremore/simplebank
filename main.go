package main

import (
	"database/sql"
	"github/aspiremore/simplebank/api"
	db "github/aspiremore/simplebank/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)
const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	local_address = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("We cannot connect to the DB:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err  = server.Start(local_address)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}

}