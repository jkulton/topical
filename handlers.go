package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

// HandlerHelper provides useful helpers to handler functions
type HandlerHelper struct {
	templates *template.Template
	storage   *Storage
}

func userFromContext(ctx context.Context) (*User, error) {
	userValue := ctx.Value(ContextUserKey)

	if userValue == nil {
		return nil, errors.New("User not found")
	}

	user := userValue.(User)

	return &user, nil
}

// TopicList renders a list of recent topics with message counts in order of most recent post
func TopicList(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())
		topics, err := h.storage.GetRecentTopics()

		if err != nil {
			log.Print("Error calling getRecentTopics")
			log.Panic(err)
		}

		payload := struct {
			Topics []Topic
			User   *User
		}{Topics: topics, User: user}

		h.templates.ExecuteTemplate(w, "list", payload)
	})
}

// TopicShow renders a topic with it's associated threaded messages
func TopicShow(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())
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

		// TODO: improve this check. helper should just return `nil` outright.
		// Redirect home with a toast message in the header.
		if topic.ID == nil {
			w.Write([]byte("404 topic not found"))
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
		user, err := userFromContext(r.Context())

		if err != nil {
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

		_, err = h.storage.CreateMessage(&message)

		if err != nil {
			log.Print("Error calling createMessage")
			log.Panic(err)
			// TODO: toast error?
			http.Redirect(w, r, "/topics", 302)
		}

		http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), 302)
	})
}

// TopicNew renders a form for creating a new topic
func TopicNew(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())
		payload := struct{ User *User }{User: user}

		h.templates.ExecuteTemplate(w, "new-topic", payload)
	})
}

// TopicCreate creates a new topic based on inputs from client
func TopicCreate(h *HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromContext(r.Context())

		if err != nil {
			http.Redirect(w, r, "/topics", 302)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		authorTheme := user.Theme
		authorInitials := user.Initials

		topic, err := h.storage.CreateTopic(title)

		if err != nil {
			log.Println(err.Error())
			log.Panic(err)
		}

		id := topic.ID

		message := Message{
			TopicID:        id,
			Content:        content,
			AuthorTheme:    authorTheme,
			AuthorInitials: authorInitials,
		}

		_, err = h.storage.CreateMessage(&message)

		if err != nil {
			log.Println(err.Error())
			log.Panic(err)
		}

		// TODO: check for best status code on creation redirect
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
		user, _ := userFromContext(r.Context())

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
			// TODO: add Flash about submission being invalid
			http.Redirect(w, r, "/join", 302)
		}

		theme, err := strconv.Atoi(r.FormValue("theme"))

		if err != nil {
			panic(err)
		}

		u := User{initials, theme}
		j, err := json.Marshal(u)

		if err != nil {
			panic(err)
		}

		encoded := base64.StdEncoding.EncodeToString([]byte(j))
		c := http.Cookie{Name: "u", Value: encoded}

		http.SetCookie(w, &c)
		http.Redirect(w, r, "/topics", 302)
	})
}
