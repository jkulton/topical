package api

import (
	"net/http"
)

// JoinShow renders the page allowing a user to log in
func (t *TopicalAPI) JoinShow(w http.ResponseWriter, r *http.Request) {
	user, _ := t.session.GetUser(r)

	// Redirect to homepage if user exists
	if user != nil {
		http.Redirect(w, r, "/topics", 302)
		return
	}

	t.templates.ExecuteTemplate(w, "join", nil)
}
