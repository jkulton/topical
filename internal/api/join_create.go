package api

import (
	"github.com/jkulton/topical/internal/models"
	"net/http"
	"regexp"
	"strconv"
)

// JoinCreate accepts a payload of user info and saves the user in a session
func (t *TopicalAPI) JoinCreate(w http.ResponseWriter, r *http.Request) {
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

	u := &models.User{initials, theme}

	if err := t.session.SaveUser(u, r, w); err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/topics", 302)
}
