package main

import (
	"database/sql"
	"github.com/jkulton/topical/internal/config"
	_ "github.com/lib/pq" // Postgres driver
	"io/ioutil"
	"log"
)

func main() {
	ac := config.ParseAppConfig()
	db, err := sql.Open("postgres", ac.DBConnectionURI)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	content, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	schema := string(content)

	_, err = db.Exec(schema)

	if err != nil {
		log.Fatal(err)
	}

}
