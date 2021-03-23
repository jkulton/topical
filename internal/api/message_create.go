package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jkulton/topical/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// MessageCreate accepts a form POST, creating a message within a given Topic
func (api *TopicalAPI) MessageCreate(w http.ResponseWriter, r *http.Request) {
	user, err := api.session.GetUser(r)

	if err != nil {
		api.session.SaveFlash("Please join to create a message", r, w)
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
		api.session.SaveFlash("Content cannot be blank", r, w)
		http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), 302)
		return
	}

	message := models.Message{
		TopicID:        &id,
		Content:        content,
		AuthorTheme:    authorTheme,
		AuthorInitials: authorInitials,
	}

	if _, err := api.storage.CreateMessage(&message); err != nil {
		// log.Print(err.Error())
		api.session.SaveFlash("Error creating message", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/topics/%d", id), 302)
}
