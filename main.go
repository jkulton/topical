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
	// Create server providing DB, templates, and cookie features to handlers
	session := sessions.NewCookieStore([]byte("69d3f5e8-d6b2-46ee-ad47-da2a12fb67ee"))
	storage, err := NewStorage("sqlite3", "./tinyboard.db")
	templates := GenerateTemplates("views/*.gohtml")
	s := &TopicServer{templates, storage, &Session{session}}

	if err != nil {
		log.Print("Error initializing storage")
		panic(err)
	}

	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/topics", s.TopicCreate).Methods("POST")
	r.HandleFunc("/topics/new", s.TopicNew).Methods("GET")
	r.HandleFunc("/topics/{id:[0-9]+}/messages", s.MessageCreate).Methods("POST")
	r.HandleFunc("/join", s.JoinShow).Methods("GET")
	r.HandleFunc("/join", s.JoinCreate).Methods("POST")
	r.HandleFunc("/", s.TopicList).Methods("GET")
	r.HandleFunc("/topics", s.TopicList).Methods("GET")
	r.HandleFunc("/topics/", s.TopicList).Methods("GET")
	r.HandleFunc("/topics/{id:[0-9]+}", s.TopicShow).Methods("GET")

	// Serve FE assets under `/static`
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Middleware Registration
	r.Use(RequestLoggerMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
