package driver

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func ConnectDB() *sql.DB {
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANT_SQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", pgURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
