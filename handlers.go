package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// HandlerHelper provides useful helpers to handler functions
type HandlerHelper struct {
	templates *template.Template
	storage   *Storage
	session   *sessions.CookieStore
}

func userFromSession(s *sessions.CookieStore, r *http.Request) (*User, error) {
	session, _ := s.Get(r, "u")
	val := session.Values["user"]
	var u *User

	if val == nil {
		return nil, errors.New("User not found")
	}

	json.Unmarshal([]byte(val.(string)), &u)
	return u, nil
}

func saveUserToSession(u *User, s *sessions.CookieStore, r *http.Request, w http.ResponseWriter) error {
	session, _ := s.Get(r, "u")
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

func saveFlashToSession(message string, s *sessions.CookieStore, r *http.Request, w http.ResponseWriter) error {
	session, err := s.Get(r, "flashes")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	session.AddFlash(message)
	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func flashesFromSession(s *sessions.CookieStore, r *http.Request, w http.ResponseWriter) ([]string, error) {
	session, _ := s.Get(r, "flashes")
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

// TopicList renders a list of recent topics with message counts in order of most recent post
func TopicList(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flashes, err := flashesFromSession(h.session, r, w)
		user, _ := userFromSession(h.session, r)
		topics, err := h.storage.GetRecentTopics()

		if err != nil {
			log.Print("Error calling getRecentTopics")
			log.Panic(err)
		}

		payload := struct {
			Topics  []Topic
			User    *User
			Flashes []string
		}{Topics: topics, User: user, Flashes: flashes}

		h.templates.ExecuteTemplate(w, "list", payload)
	})
}

// TopicShow renders a topic with it's associated threaded messages
func TopicShow(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromSession(h.session, r)
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Print("Error parsing route id")
			log.Panic(err)
		}

		topic, err := h.storage.GetTopic(id)

		if err != nil {
			log.Print("Error calling getTopic")
			log.Panic(err)
		}

		if topic.ID == nil {
			if err := saveFlashToSession("Topic not found", h.session, r, w); err != nil {
				panic(err)
			}
			http.Redirect(w, r, "/topics", 302)
			return
		}

		payload := struct {
			Topic *Topic
			User  *User
		}{topic, user}

		h.templates.ExecuteTemplate(w, "show", payload)
	})
}

// MessageCreate accepts a form POST, creating a message within a given Topic
func MessageCreate(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromSession(h.session, r)

		if err != nil {
			if err := saveFlashToSession("Please join to create a message", h.session, r, w); err != nil {
				panic(err)
			}
			http.Redirect(w, r, "/topics", 302)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Print("Error parsing route id")
			log.Panic(err)
		}

		content := r.FormValue("content")
		authorTheme := user.Theme
		authorInitials := user.Initials

		message := Message{
			TopicID:        &id,
			Content:        content,
			AuthorTheme:    authorTheme,
			AuthorInitials: authorInitials,
		}

		if _, err := h.storage.CreateMessage(&message); err != nil {
			log.Panic(err)
			if err := saveFlashToSession("Error creating message", h.session, r, w); err != nil {
				panic(err)
			}
			http.Redirect(w, r, "/topics", 302)
		}

		http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), 302)
	})
}

// TopicNew renders a form for creating a new topic
func TopicNew(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromSession(h.session, r)

		if err != nil {
			log.Print("User does not exist, redirecting home")
			if err := saveFlashToSession("Log in to post a message", h.session, r, w); err != nil {
				panic(err)
			}
			http.Redirect(w, r, "/topics", 302)
			return
		}

		payload := struct{ User *User }{User: user}
		h.templates.ExecuteTemplate(w, "new-topic", payload)
	})
}

// TopicCreate creates a new topic based on inputs from client
func TopicCreate(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromSession(h.session, r)

		if err != nil {
			log.Print("User does not exist, redirecting home")
			if err := saveFlashToSession("Log in to post a topic", h.session, r, w); err != nil {
				panic(err)
			}
			http.Redirect(w, r, "/topics", 302)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		topic, err := h.storage.CreateTopic(title)

		if err != nil {
			log.Println(err.Error())
			log.Panic(err)
		}

		message := Message{
			TopicID:        topic.ID,
			Content:        content,
			AuthorTheme:    user.Theme,
			AuthorInitials: user.Initials,
		}

		_, err = h.storage.CreateMessage(&message)

		if err != nil {
			log.Println(err.Error())
			log.Panic(err)
		}

		http.Redirect(w, r, fmt.Sprintf("/topics/%d", *topic.ID), 302)
	})
}

// SettingsUpdate is a POST request used for setting the general settings for the user
func SettingsUpdate(h HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Settings struct {
			Theme string
		}

		theme := r.FormValue("theme")

		if theme != "dark" {
			theme = "light"
		}

		s := Settings{theme}
		b, err := json.Marshal(s)

		if err != nil {
			panic(err)
		}

		c := http.Cookie{
			Name:  "s",
			Value: string(b),
		}
		http.SetCookie(w, &c)

		http.Redirect(w, r, "/topics", 302)
	})
}

// JoinShow renders the page allowing a user to 'join', which really just creates a
// cookie in their local browser with user information.
func JoinShow(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromSession(h.session, r)

		// Redirect to homepage if user exists
		if user != nil {
			http.Redirect(w, r, "/topics", 302)
			return
		}

		h.templates.ExecuteTemplate(w, "join", nil)
	})
}

// JoinCreate uses input from the POSTed form to store the user's
// information in a cookie.
func JoinCreate(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initials := r.FormValue("initials")
		matched, err := regexp.Match("^[A-Z]{2}$", []byte("JK"))

		if err != nil {
			panic(err)
		}

		if matched == false {
			http.Redirect(w, r, "/join", 302)
			return
		}

		theme, err := strconv.Atoi(r.FormValue("theme"))

		if err != nil {
			panic(err)
		}

		u := &User{initials, theme}

		if err := saveUserToSession(u, h.session, r, w); err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/topics", 302)
	})
}
