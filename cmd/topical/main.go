package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jkulton/topical/internal/api"
	"github.com/jkulton/topical/internal/config"
	"github.com/jkulton/topical/internal/middleware"
	"github.com/jkulton/topical/internal/session"
	"github.com/jkulton/topical/internal/storage"
	"github.com/jkulton/topical/internal/templates"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Grab configuration from flags or ENV
	ac := config.ParseAppConfig()

	// DB Setup
	db, err := sql.Open("sqlite3", ac.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize session, HTML templates, and storage interface
	session := session.NewSession(ac.SessionKey)
	templates, err := templates.GenerateTemplates("./web/views/*.gohtml")

	if err != nil {
		log.Fatal(err)
	}

	storage := storage.New(db)

	// Create API & router, register routes
	a := api.New(templates, storage, session)
	r := mux.NewRouter()
	a.RegisterRoutes(r)

	// Serve FE assets under `/static`
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// Middleware Registration
	r.Use(middleware.RequestLogger)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ac.Port), r))
}
