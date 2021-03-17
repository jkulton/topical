package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type TopicalSession interface {
	GetUser(r *http.Request) (*User, error)
	SaveUser(u *User, r *http.Request, w http.ResponseWriter) error
	SaveFlash(message string, r *http.Request, w http.ResponseWriter) error
	GetFlashes(r *http.Request, w http.ResponseWriter) ([]string, error)
}

// Session is a struct which wraps a gorilla/sessions CookieStore
// and implements a few methods used for retrieval and storage of session items.
type Session struct {
	session *sessions.CookieStore
}

// User is a struct representing a user account, the user stored
// in a simple cookie and defines a name and theme for messages.
type User struct {
	Initials string
	Theme    int
}

// GetUser returns the the User from the session, if present
func (s *Session) GetUser(r *http.Request) (*User, error) {
	session, _ := s.session.Get(r, "u")
	val := session.Values["user"]
	var u *User

	if val == nil {
		return nil, errors.New("User not found")
	}

	json.Unmarshal([]byte(val.(string)), &u)
	return u, nil
}

// SaveUser saves a user to the session
func (s *Session) SaveUser(u *User, r *http.Request, w http.ResponseWriter) error {
	session, _ := s.session.Get(r, "u")
	j, err := json.Marshal(u)

	if err != nil {
		return errors.New("Unable to save user")
	}

	session.Values["user"] = string(j)

	if err := session.Save(r, w); err != nil {
		return errors.New("Unable to save user")
	}

	return nil
}

// SaveFlash saves a flash message to the session
func (s *Session) SaveFlash(message string, r *http.Request, w http.ResponseWriter) error {
	session, err := s.session.Get(r, "flashes")

	if err != nil {
		return err
	}

	session.AddFlash(message)
	err = session.Save(r, w)

	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

// GetFlashes returns all flash messages stored in the session.
// Note that the way flash messages work they are deleted after being retrieved.
func (s *Session) GetFlashes(r *http.Request, w http.ResponseWriter) ([]string, error) {
	session, _ := s.session.Get(r, "flashes")
	flashStrings := []string{}
	flashes := session.Flashes()

	if len(flashes) == 0 {
		return nil, errors.New("No flashes found")
	}

	for _, flash := range flashes {
		flashStrings = append(flashStrings, flash.(string))
	}

	session.Save(r, w)

	return flashStrings, nil
}
