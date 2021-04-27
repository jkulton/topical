package api

import (
	"html/template"

	"github.com/gorilla/mux"
	"github.com/jkulton/topical/internal/session"
	"github.com/jkulton/topical/internal/storage"
)

// TopicalAPI represents an API instance, with internal state for
// templates, storage, and session used by handlers.
type TopicalAPI struct {
	templates *template.Template
	storage   storage.TopicalStore
	session   session.TopicalSession
}

// New returns a new TopicalAPI instance
func New(templates *template.Template, storage storage.TopicalStore, session session.TopicalSession) *TopicalAPI {
	return &TopicalAPI{templates, storage, session}
}

// RegisterRoutes registers handler functions defined in this package on a router instance
func (t *TopicalAPI) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/topics", t.TopicCreate).Methods("POST")
	r.HandleFunc("/topics/new", t.TopicNew).Methods("GET")
	r.HandleFunc("/topics/{id:[0-9]+}/messages", t.MessageCreate).Methods("POST")
	r.HandleFunc("/join", t.JoinShow).Methods("GET")
	r.HandleFunc("/join", t.JoinCreate).Methods("POST")
	r.HandleFunc("/", t.TopicList).Methods("GET")
	r.HandleFunc("/topics", t.TopicList).Methods("GET")
	r.HandleFunc("/topics/", t.TopicList).Methods("GET")
	r.HandleFunc("/topics/{id:[0-9]+}", t.TopicShow).Methods("GET")
}
