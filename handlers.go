package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// TopicServer provides HTTP handlers to routes, and template and storage helpers
// to the HTTP handlers themselves.
type TopicServer struct {
	templates *template.Template
	storage   TopicalStore
	session   TopicalSession
}

// TopicList renders a list of recent topics with message counts in order of most recent post
func (s *TopicServer) TopicList(w http.ResponseWriter, r *http.Request) {
	flashes, err := s.session.GetFlashes(r, w)
	user, _ := s.session.GetUser(r)
	topics, err := s.storage.GetRecentTopics()

	if err != nil {
		log.Print("Error calling getRecentTopics")
		log.Panic(err)
	}

	payload := struct {
		Topics  []Topic
		User    *User
		Flashes []string
	}{Topics: topics, User: user, Flashes: flashes}

	s.templates.ExecuteTemplate(w, "list", payload)
}

// TopicShow renders a topic with it's associated threaded messages
func (s *TopicServer) TopicShow(w http.ResponseWriter, r *http.Request) {
	flashes, _ := s.session.GetFlashes(r, w)
	user, _ := s.session.GetUser(r)
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		log.Print("Error parsing route id")
		log.Panic(err)
	}

	topic, err := s.storage.GetTopic(id)

	if err != nil {
		log.Print("Error calling getTopic")
		log.Panic(err)
	}

	if topic.ID == nil {
		s.session.SaveFlash("Topic not found", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	payload := struct {
		Topic   *Topic
		User    *User
		Flashes []string
	}{topic, user, flashes}

	s.templates.ExecuteTemplate(w, "show", payload)
}

// MessageCreate accepts a form POST, creating a message within a given Topic
func (s *TopicServer) MessageCreate(w http.ResponseWriter, r *http.Request) {
	user, err := s.session.GetUser(r)

	if err != nil {
		s.session.SaveFlash("Please join to create a message", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Print("Error parsing route id")
		log.Panic(err)
	}

	content := strings.TrimSpace(r.FormValue("content"))
	authorTheme := user.Theme
	authorInitials := user.Initials

	if content == "" {
		s.session.SaveFlash("Content cannot be blank", r, w)
		http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), 302)
		return
	}

	message := Message{
		TopicID:        &id,
		Content:        content,
		AuthorTheme:    authorTheme,
		AuthorInitials: authorInitials,
	}

	if _, err := s.storage.CreateMessage(&message); err != nil {
		// log.Print(err.Error())
		s.session.SaveFlash("Error creating message", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), 302)
}

// TopicNew renders a form for creating a new topic
func (s *TopicServer) TopicNew(w http.ResponseWriter, r *http.Request) {
	flashes, _ := s.session.GetFlashes(r, w)
	user, err := s.session.GetUser(r)

	if err != nil {
		// log.Print("User does not exist, redirecting home")
		s.session.SaveFlash("Log in to post a message", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	payload := struct {
		User    *User
		Flashes []string
	}{user, flashes}

	s.templates.ExecuteTemplate(w, "new-topic", payload)
}

// TopicCreate creates a new topic based on inputs from client
func (s *TopicServer) TopicCreate(w http.ResponseWriter, r *http.Request) {
	user, err := s.session.GetUser(r)

	if err != nil {
		// log.Print("User does not exist, redirecting home")
		s.session.SaveFlash("Log in to post a topic", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	if content == "" || title == "" {
		s.session.SaveFlash("Inputs cannot be blank", r, w)
		http.Redirect(w, r, "/topics/new", 302)
		return
	}

	topic, err := s.storage.CreateTopic(title)

	if err != nil {
		log.Println("CreateTopic error")
		log.Println(err.Error())
	}

	message := Message{
		TopicID:        topic.ID,
		Content:        content,
		AuthorTheme:    user.Theme,
		AuthorInitials: user.Initials,
	}

	_, err = s.storage.CreateMessage(&message)

	if err != nil {
		log.Println(err.Error())
		log.Panic(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", *topic.ID), 302)
}

// JoinShow renders the page allowing a user to log in
func (s *TopicServer) JoinShow(w http.ResponseWriter, r *http.Request) {
	user, _ := s.session.GetUser(r)

	// Redirect to homepage if user exists
	if user != nil {
		http.Redirect(w, r, "/topics", 302)
		return
	}

	s.templates.ExecuteTemplate(w, "join", nil)
}

// JoinCreate accepts a payload of user info and saves the user in a session
func (s *TopicServer) JoinCreate(w http.ResponseWriter, r *http.Request) {
	initials := r.FormValue("initials")
	matched, err := regexp.Match("^[A-Z]{2}$", []byte(initials))

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

	if err := s.session.SaveUser(u, r, w); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/topics", 302)
}
