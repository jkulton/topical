package main

import (
	"database/sql"
	"log"
	"os"

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

	content, err := os.ReadFile("./seeds.sql")
	if err != nil {
		log.Fatal(err)
	}

	schema := string(content)
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}
}
