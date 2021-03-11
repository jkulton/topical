package main

/**

Remaining TODOs:

- update to gorilla/sessions (for Flash messages and a little security around user)
- settings PUT for dark/light mode (use gorilla/sessions)
- update to Postgres for Heroku deployment
- create some logging Middleware?
- UI redesign

DONE:

- refactor UserMiddleware
- redirect home on POST endpoints when user wasn't parsed
- make sure we have validation around user initials being two characters
- break app into multiple files

**/

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Set up helper providing DB interface and HTML templates to route handlers
	templates := GenerateTemplates("views/*.gohtml")
	storage, err := NewStorage("sqlite3", "./tinyboard.db")

	if err != nil {
		log.Print("Error initializing storage")
		panic(err)
	}

	defer storage.db.Close()
	h := &HandlerHelper{templates, storage}
	r := mux.NewRouter()

	// Serve files under `/static` (FE assets)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Register routes
	r.HandleFunc("/topics", TopicCreate(h)).Methods("POST").Name("TopicCreate")
	r.HandleFunc("/topics/new", TopicNew(h)).Methods("GET").Name("TopicNew")
	r.HandleFunc("/topics/{id}/messages", MessageCreate(h)).Methods("POST").Name("MessageCreate")
	r.HandleFunc("/join", JoinShow(h)).Methods("GET").Name("JoinShow")
	r.HandleFunc("/join", JoinCreate(h)).Methods("POST").Name("JoinCreate")
	r.HandleFunc("/", TopicList(h)).Methods("GET").Name("TopicList")
	r.HandleFunc("/topics", TopicList(h)).Methods("GET")
	r.HandleFunc("/topics/", TopicList(h)).Methods("GET")
	r.HandleFunc("/topics/{id:[0-9]+}", TopicShow(h)).Methods("GET").Name("TopicShow")

	// Middleware Registration
	r.Use(RequestLoggerMiddleware)
	r.Use(UserMiddleware)

	// Protected Routes are those which require a user be logged in.
	protectedRoutes := []string{"TopicCreate", "TopicNew", "MessageCreate"}
	r.Use(ProtectedRouteMiddleware(protectedRoutes))

	log.Fatal(http.ListenAndServe(":8080", r))
}
