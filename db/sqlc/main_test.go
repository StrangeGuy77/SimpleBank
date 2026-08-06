package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "$(DB_SOURCE)"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Couldn't connect to the database: ", err)
	}

	testDB.SetMaxIdleConns(10)
	testDB.SetMaxOpenConns(10)
	testQueries = New(testDB)

	os.Exit(m.Run())
}
