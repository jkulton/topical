package api

import (
	"github.com/jkulton/topical/internal/models"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// JoinCreate accepts a payload of user info and saves the user in a session
func (t *TopicalAPI) JoinCreate(w http.ResponseWriter, r *http.Request) {
	initials := strings.ToUpper(r.FormValue("initials"))
	matched, err := regexp.Match("^[A-Z]{2}$", []byte(initials))

	if err != nil {
		log.Print("Error creating user")
		log.Print(err.Error())
		http.Redirect(w, r, "/join", 302)
	}

	if matched == false {
		http.Redirect(w, r, "/join", 302)
		return
	}

	theme, err := strconv.Atoi(r.FormValue("theme"))

	if err != nil {
		log.Print("Error creating user")
		log.Print(err.Error())
		http.Redirect(w, r, "/join", 302)
		return
	}

	u := &models.User{Initials: initials, Theme: theme}

	if err := t.session.SaveUser(u, r, w); err != nil {
		log.Print("Error creating user")
		log.Print(err.Error())
		http.Redirect(w, r, "/join", 302)
		return
	}

	http.Redirect(w, r, "/topics", 302)
}
