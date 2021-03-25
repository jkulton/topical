package api

import (
	"github.com/jkulton/topical/internal/models"
	"log"
	"net/http"
)

// TopicList renders a list of recent topics with message counts in order of most recent post
func (api *TopicalAPI) TopicList(w http.ResponseWriter, r *http.Request) {
	flashes, err := api.session.GetFlashes(r, w)
	user, _ := api.session.GetUser(r)
	topics, err := api.storage.GetRecentTopics()

	if err != nil {
		log.Print("Error getting recent topics")
		log.Print(err.Error())
		api.templates.ExecuteTemplate(w, "error", nil)
		return
	}

	payload := struct {
		Topics  []models.Topic
		User    *models.User
		Flashes []string
	}{Topics: topics, User: user, Flashes: flashes}

	api.templates.ExecuteTemplate(w, "list", payload)
}
