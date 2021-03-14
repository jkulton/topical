package main

/**

Remaining TODOs:

- settings PUT for dark/light mode (session data)
- update to Postgres for Heroku deployment
- TESTS TESTS TESTS
- implement new UI

DONE:

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

**/

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func main() {
	// Create helper providing DB, templates, and cookie features to handlers
	session := sessions.NewCookieStore([]byte("69d3f5e8-d6b2-46ee-ad47-da2a12fb67ee"))
	storage, err := NewStorage("sqlite3", "./tinyboard.db")
	templates := GenerateTemplates("views/*.gohtml")

	if err != nil {
		log.Print("Error initializing storage")
		panic(err)
	}

	h := &HandlerHelper{templates, storage, session}

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/topics", TopicCreate(h)).Methods("POST").Name("TopicCreate")
	r.HandleFunc("/topics/new", TopicNew(h)).Methods("GET").Name("TopicNew")
	r.HandleFunc("/topics/{id}/messages", MessageCreate(h)).Methods("POST").Name("MessageCreate")
	r.HandleFunc("/join", JoinShow(h)).Methods("GET").Name("JoinShow")
	r.HandleFunc("/join", JoinCreate(h)).Methods("POST").Name("JoinCreate")
	r.HandleFunc("/", TopicList(h)).Methods("GET").Name("TopicList")
	r.HandleFunc("/topics", TopicList(h)).Methods("GET")
	r.HandleFunc("/topics/", TopicList(h)).Methods("GET")
	r.HandleFunc("/topics/{id:[0-9]+}", TopicShow(h)).Methods("GET").Name("TopicShow")

	// Serve FE assets under `/static`
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Middleware Registration
	r.Use(RequestLoggerMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
