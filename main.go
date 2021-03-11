package main

/**

Remaining TODOs:

- update to gorilla/sessions (for Flash messages and a little security around user)
- settings PUT for dark/light mode (use gorilla/sessions)
- update to Postgres for Heroku deployment
- break app into multiple files
- create some logging Middleware?
- UI redesign

DONE:

- refactor UserMiddleware
- redirect home on POST endpoints when user wasn't parsed
- make sure we have validation around user initials being two characters

**/

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// Message represents a message entity, messages have a M:1 relationship with Topics
type Message struct {
	TopicID        *int
	Content        string
	AuthorInitials string
	Posted         string
	AuthorTheme    int
}

// Topic represents a topic entity, one Topic can contain many Messages
type Topic struct {
	ID             *int
	Title          string
	Messages       *[]Message
	MessageCount   *int
	AuthorInitials *string
	AuthorTheme    *string
}

// HandlerHelper is a struct designed for injecting helpers instances into a handler function
type HandlerHelper struct {
	db        string
	templates *template.Template
	user      *User
}

// User is a struct representing a user account, the user stored
// in a simple cookie and defines a name and theme for messages.
type User struct {
	Initials string
	Theme    int
}

// Route represents an HTTP endpoint for the application
type Route struct {
	method  string
	path    string
	handler http.HandlerFunc
	name    string
}

type ContextKey string

const ContextUserKey ContextKey = "user"

// ProtectedRouteMiddleware redirects home if a protected route is attempted
// to be accessed without a user present in Context.
func ProtectedRouteMiddleware(protectedRouteNames []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isProtectedRoute := false
			routeName := mux.CurrentRoute(r).GetName()
			_, err := userFromContext(r.Context())

			for _, protected := range protectedRouteNames {
				if routeName == protected {
					isProtectedRoute = true
				}
			}

			if err != nil && isProtectedRoute {
				log.Print("User not present on protected route, redirecting home")
				http.Redirect(w, r, "/topics", 302)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UserMiddleware extracts the user from request cookie and stores user in context
func UserMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var user User
			cookie, err := r.Cookie("u")

			if err != nil {
				ctx = context.WithValue(ctx, ContextUserKey, nil)
			} else {
				cookieByte, _ := base64.StdEncoding.DecodeString(cookie.Value)
				cookieStr := string(cookieByte)
				json.Unmarshal([]byte(cookieStr), &user)
				ctx = context.WithValue(ctx, ContextUserKey, user)
			}

			rWithUser := r.WithContext(ctx)
			next.ServeHTTP(w, rWithUser)
		})
	}
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
func TopicList(helper HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())
		topics, err := getRecentTopics()

		if err != nil {
			log.Print("Error calling getRecentTopics")
			log.Panic(err)
		}

		payload := struct {
			Topics []Topic
			User   *User
		}{Topics: topics, User: user}

		helper.templates.ExecuteTemplate(w, "list", payload)
	})
}

// TopicShow renders a topic with it's associated threaded messages
func TopicShow(helper HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Print("Error parsing route id")
			log.Panic(err)
		}

		topic, err := getTopic(id)

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

		helper.templates.ExecuteTemplate(w, "show", payload)
	})
}

// MessageCreate accepts a form POST, creating a message within a given Topic
func MessageCreate(helper HandlerHelper) http.HandlerFunc {
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

		_, err = createMessage(&message)

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
func TopicNew(helper HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())
		payload := struct{ User *User }{User: user}

		helper.templates.ExecuteTemplate(w, "new-topic", payload)
	})
}

// TopicCreate creates a new topic based on inputs from client
func TopicCreate(helper HandlerHelper) http.HandlerFunc {
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

		topic, err := createTopic(title)

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

		_, err = createMessage(&message)

		if err != nil {
			log.Println(err.Error())
			log.Panic(err)
		}

		// TODO: check for best status code on creation redirect
		http.Redirect(w, r, fmt.Sprintf("/topics/%d", *topic.ID), 302)
	})
}

// SettingsUpdate is a POST request used for setting the general settings for the user
func SettingsUpdate(helper HandlerHelper) http.HandlerFunc {
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
func JoinShow(helper HandlerHelper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := userFromContext(r.Context())

		// Redirect to homepage if user exists
		if user != nil {
			http.Redirect(w, r, "/topics", 302)
			return
		}

		helper.templates.ExecuteTemplate(w, "join", nil)
	})
}

// JoinCreate uses input from the POSTed form to store the user's
// information in a cookie.
func JoinCreate(helper HandlerHelper) http.HandlerFunc {
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

// createHelper returns a HandlerHelper instance for use in helpers
func createHelper() HandlerHelper {
	funcMap := template.FuncMap{
		"noescape": func(str string) template.HTML { return template.HTML(str) },
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob("views/*.gohtml")

	if err != nil {
		log.Print("Error compiling templates:")
		log.Print(err)
		panic(err)
	}

	return HandlerHelper{
		db:        "foo",
		templates: templates,
	}
}

func main() {
	helper := createHelper()

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	routes := []Route{
		{"POST", "/topics", TopicCreate(helper), "TopicCreate"},
		{"GET", "/topics/new", TopicNew(helper), "TopicNew"},
		{"POST", "/topics/{id}/messages", MessageCreate(helper), "MessageCreate"},
		{"GET", "/join", JoinShow(helper), "JoinShow"},
		{"POST", "/join", JoinCreate(helper), "JoinCreate"},
		{"GET", "/", TopicList(helper), "TopicList1"},
		{"GET", "/topics", TopicList(helper), "TopicList2"},
		{"GET", "/topics/", TopicList(helper), "TopicList3"},
		{"GET", "/topics/{id:[0-9]+}", TopicShow(helper), "TopicShow"},
		// {"POST", "/settings", SettingsUpdate(helper), "SettingsUpdate"}
	}

	protectedRouteNames := []string{"TopicCreate", "TopicNew", "MessageCreate"}

	for _, route := range routes {
		r.HandleFunc(route.path, route.handler).Methods(route.method).Name(route.name)
	}

	r.Use(UserMiddleware())
	r.Use(ProtectedRouteMiddleware(protectedRouteNames))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func dbConn(driver, database string) (db *sql.DB, error error) {
	db, err := sql.Open(driver, database)
	if err != nil {
		return nil, error
	}
	return db, nil
}

func getTopic(id int) (*Topic, error) {
	topic := Topic{}
	messages := []Message{}
	query := `
		SELECT
			topics.id,
			topics.title,
			messages.content,
			messages.author_initials,
			messages.author_theme,
			messages.posted
		FROM topics
		INNER JOIN messages ON messages.topic_id = topics.id
		WHERE topics.id = ?
		ORDER BY posted;
	`

	db, err := dbConn("sqlite3", "./tinyboard.db")

	// TODO: DRY all this out if possible
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query, id)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var content string
		var authorInitials string
		var posted string
		var authorTheme int

		err = rows.Scan(&id, &title, &content, &authorInitials, &authorTheme, &posted)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		topic.ID = &id
		topic.Title = title

		unsafeHTML := blackfriday.MarkdownBasic([]byte(content))
		safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)

		messages = append(messages, Message{
			Content:        string(safeHTML),
			AuthorInitials: authorInitials,
			Posted:         posted,
			AuthorTheme:    authorTheme,
		})
	}
	err = rows.Err()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	topic.Messages = &messages

	return &topic, nil
}

func getRecentTopics() ([]Topic, error) {
	topics := []Topic{}
	/**
		TODO: Fix subquery for author_theme, need to join color table.
		TODO: This entire query can be made more performant.
	**/
	query := `
		SELECT DISTINCT topics.*,
			(SELECT COUNT(messages.id) FROM messages WHERE topic_id = topics.id) AS "message_count",
			(SELECT author_initials FROM messages WHERE topic_id = topics.id ORDER BY posted LIMIT 1) AS "author_initials",
			(SELECT author_theme FROM messages WHERE topic_id = topics.id ORDER BY posted LIMIT 1) AS "author_theme",
			(SELECT posted FROM messages WHERE topic_id = topics.id ORDER BY posted DESC LIMIT 1) AS "last_message"
		FROM topics
		INNER JOIN messages
		ON topics.id = messages.topic_id
		ORDER BY last_message DESC
		LIMIT 50;
	`

	db, err := sql.Open("sqlite3", "./tinyboard.db")

	// TODO: DRY all this out if possible
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id, messageCount int
		var title, authorInitials, authorTheme, posted string
		err = rows.Scan(&id, &title, &messageCount, &authorInitials, &authorTheme, &posted)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		topics = append(topics, Topic{
			ID:             &id,
			Title:          title,
			MessageCount:   &messageCount,
			AuthorInitials: &authorInitials,
			AuthorTheme:    &authorTheme,
		})
	}
	err = rows.Err()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return topics, nil
}

func createMessage(m *Message) (*Message, error) {
	db, err := dbConn("sqlite3", "./tinyboard.db")

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	defer db.Close()

	sql := `INSERT INTO messages (topic_id, content, author_initials, author_theme) VALUES (?, ?, ?, ?)`
	query, err := db.Prepare(sql)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	_, err = query.Exec(&m.TopicID, m.Content, m.AuthorInitials, m.AuthorTheme)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return m, nil
}

func createTopic(title string) (*Topic, error) {
	db, err := dbConn("sqlite3", "./tinyboard.db")

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	defer db.Close()

	// SQLite doesn't support INSERT ... RETURNING, this is a workaround for that.
	query, err := db.Prepare(`INSERT INTO topics (title) VALUES (?)`)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	res, err := query.Exec(title)

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	idInt64, err := res.LastInsertId()

	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	id := int(idInt64)

	return &Topic{ID: &id, Title: title}, nil
}
