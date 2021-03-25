package main

import (
	"database/sql"
	"github.com/jkulton/topical/internal/config"
	_ "github.com/lib/pq" // Postgres driver
	"log"
)

func main() {
	ac := config.ParseAppConfig()
	db, err := sql.Open("postgres", ac.DBConnectionURI)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`DROP TABLE messages`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`DROP TABLE topics`)

	if err != nil {
		panic(err)
	}

}
