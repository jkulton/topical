package main

/**

- update to gorilla/sessions for cookie stuff
- refactor UserMiddleware
- redirect home on POST endpoints when user wasn't parsed
- make sure we have validation around user initials being two characters
- break app into multiple files
- create some logging Middleware?
- get rid of ProtectedRouteMiddleware? we still do a second check in the handler anyhow
	- only three routes are protected so seems safer to just do it manually for now
- add Flash message errors
- UI redesign
- TESTS TESTS TESTS
- settings PUT for dark/light mode (session data)
- implement new UI
- update to Postgres for Heroku deployment
- pass config (ports, ssl, db) values through env
- move main.go code into /cmd for bin
- dark mode option?

**/

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jkulton/board/internal/api"
	"github.com/jkulton/board/internal/config"
	"github.com/jkulton/board/internal/middleware"
	"github.com/jkulton/board/internal/session"
	"github.com/jkulton/board/internal/storage"
	"github.com/jkulton/board/internal/templates"
	_ "github.com/lib/pq" // Postgres driver
	"log"
	"net/http"
)

func main() {
	// Grab configuration from flags or ENV
	ac := config.ParseAppConfig()

	// DB Setup
	db, err := sql.Open("postgres", ac.DBConnectionURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize session, HTML templates, and storage interface
	session := session.NewSession(ac.SessionKey)
	templates := templates.GenerateTemplates("./web/views/*.gohtml")
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
