package api

import (
	"github.com/jkulton/board/internal/models"
	"net/http"
)

// TopicNew renders a form for creating a new topic
func (api *TopicalAPI) TopicNew(w http.ResponseWriter, r *http.Request) {
	flashes, _ := api.session.GetFlashes(r, w)
	user, err := api.session.GetUser(r)

	if err != nil {
		// log.Print("User does not exist, redirecting home")
		api.session.SaveFlash("Log in to post a message", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	payload := struct {
		User    *models.User
		Flashes []string
	}{user, flashes}

	api.templates.ExecuteTemplate(w, "new-topic", payload)
}
