package main

import (
	"database/sql"
	"github.com/aspiremore/simplebank/api"
	db "github.com/aspiremore/simplebank/db/sqlc"
	"github.com/aspiremore/simplebank/db/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to load configurations", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("We cannot connect to the DB:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err  = server.Start(config.Server_address)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}

}