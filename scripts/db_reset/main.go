package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/jkulton/topical/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ac := config.ParseAppConfig()

	db, err := sql.Open("sqlite3", ac.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Drop Tables
	_, err = db.Exec(`DROP TABLE IF EXISTS messages`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS topics`)
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
