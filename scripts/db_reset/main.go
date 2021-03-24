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

	// Drop Tables
	_, err = db.Exec(`DROP TABLE messages`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`DROP TABLE topics`)

	if err != nil {
		panic(err)
	}

	// Create Tables

	content, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	schema := string(content)

	_, err = db.Exec(schema)

	if err != nil {
		log.Fatal(err)
	}

	// Seed Tables
	content, err = ioutil.ReadFile("./seeds.sql")
	if err != nil {
		log.Fatal(err)
	}

	schema = string(content)

	_, err = db.Exec(schema)

	if err != nil {
		log.Fatal(err)
	}
}
