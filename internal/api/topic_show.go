package api

import (
	"github.com/gorilla/mux"
	"github.com/jkulton/board/internal/models"
	"log"
	"net/http"
	"strconv"
)

// TopicShow renders a topic with it's associated threaded messages
func (api *TopicalAPI) TopicShow(w http.ResponseWriter, r *http.Request) {
	flashes, _ := api.session.GetFlashes(r, w)
	user, _ := api.session.GetUser(r)
	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		log.Print("Error parsing route id")
		log.Panic(err)
	}

	topic, err := api.storage.GetTopic(id)

	if err != nil {
		log.Print("Error calling getTopic")
		log.Panic(err)
	}

	if topic.ID == nil {
		api.session.SaveFlash("Topic not found", r, w)
		http.Redirect(w, r, "/topics", 302)
		return
	}

	payload := struct {
		Topic   *models.Topic
		User    *models.User
		Flashes []string
	}{topic, user, flashes}

	api.templates.ExecuteTemplate(w, "show", payload)
}
