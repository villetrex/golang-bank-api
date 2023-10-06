package tutorial

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // add _ (blank identifier) before package name if we don't directly call any function from this library else go formatter will remove it. IN this case, we simply use the lib/pq as a driver for sqlc
	"tutorial.sqlc.dev/app/util"
)

var testQueries *Queries
var testDb *sql.DB

// the TestMain is the entry point for all unit tests
const (
	dbDriver = "postgres"
	dbSource = "postgresql://username:secretkey@localhost:5432/database_name?sslmode=disable"
)

func TestMain(t *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	var err error
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database")
	}
	testQueries = New(testDb)
	os.Exit(m.run())
}
