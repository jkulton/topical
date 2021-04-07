package api

import (
	"fmt"
	"github.com/jkulton/topical/internal/models"
	"log"
	"net/http"
	"strings"
)

// TopicCreate creates a new topic based on inputs from client
func (api *TopicalAPI) TopicCreate(w http.ResponseWriter, r *http.Request) {
	user, err := api.session.GetUser(r)

	if err != nil {
		api.session.SaveFlash("Log in to post a topic", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	if content == "" || title == "" {
		api.session.SaveFlash("Inputs cannot be blank", r, w)
		http.Redirect(w, r, "/topics/new", 302)
		return
	}

	topic, err := api.storage.CreateTopic(title)

	if err != nil {
		log.Print("Error creating topic", err.Error())
		api.templates.ExecuteTemplate(w, "error", nil)
		return
	}

	message := models.Message{
		TopicID:        topic.ID,
		Content:        content,
		AuthorTheme:    user.Theme,
		AuthorInitials: user.Initials,
	}

	_, err = api.storage.CreateMessage(&message)

	if err != nil {
		log.Print("Error creating message", err.Error())
		api.templates.ExecuteTemplate(w, "error", nil)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", *topic.ID), 302)
}
