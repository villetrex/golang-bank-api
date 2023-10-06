package main

import (
	"database/sql"
	"log"

	"tutorial.sqlc.dev/app/api"
	"tutorial.sqlc.dev/app/util"

	// add below else our coe won't be able to talk to the database
	_ "github.com/lib/pq" // add _ (blank identifier) before package name if we don't directly call any function from this library else go formatter will remove it. IN this case, we simply use the lib/pq as a driver for sqlc
)

// the TestMain is the entry point for all unit tests
// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable"
// 	serverAddress = "0.0.0.0:8080"
// )

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(" cannot load configs", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database")
	}
	store := db.NewStore(conn)
	server :=
		api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
