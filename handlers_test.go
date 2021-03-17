package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockStorage struct {
	GetTopicFunc        func(id int) (*Topic, error)
	GetRecentTopicsFunc func() ([]Topic, error)
	CreateMessageFunc   func(m *Message) (*Message, error)
	CreateTopicFunc     func(title string) (*Topic, error)
}

func (s *MockStorage) GetTopic(id int) (*Topic, error) {
	return s.GetTopicFunc(id)
}

func (s *MockStorage) GetRecentTopics() ([]Topic, error) {
	return s.GetRecentTopicsFunc()
}

func (s *MockStorage) CreateMessage(m *Message) (*Message, error) {
	return s.CreateMessageFunc(m)
}

func (s *MockStorage) CreateTopic(title string) (*Topic, error) {
	return s.CreateTopicFunc(title)
}

var (
	session   *sessions.CookieStore
	templates *template.Template
	storage   MockStorage
	s         TopicServer
)

func setupTests() {
	session = sessions.NewCookieStore([]byte("test"))
	templates = GenerateTemplates("views/*.gohtml")
	storage = MockStorage{
		GetTopicFunc: func(id int) (*Topic, error) {
			return &Topic{ID: &id, Title: "First Title"}, nil
		},
		GetRecentTopicsFunc: func() ([]Topic, error) {
			return []Topic{}, nil
		},
		CreateMessageFunc: func(m *Message) (*Message, error) {
			return nil, nil
		},
		CreateTopicFunc: func(title string) (*Topic, error) {
			return nil, nil
		},
	}

	s = TopicServer{templates, &storage, &Session{session}}
}

func assertRedirect(location string, t *testing.T, res *httptest.ResponseRecorder) {
	redirect := res.Header()["Location"][0]

	if res.Code != http.StatusFound {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusFound)
	}

	if redirect != location {
		t.Errorf("got redirect %s but wanted %s", redirect, location)
	}
}

func TestTopicShow(t *testing.T) {
	t.Run("redirects home if topic not found", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()
		vars := map[string]string{"id": "12"}
		req = mux.SetURLVars(req, vars)

		storage.GetTopicFunc = func(id int) (*Topic, error) {
			return &Topic{}, nil
		}

		s.TopicShow(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("renders topic successfully", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()
		vars := map[string]string{"id": "12"}
		req = mux.SetURLVars(req, vars)

		s.TopicShow(res, req)

		if strings.Contains(res.Body.String(), "<h2>First Title</h2>") == false {
			t.Error("response body should include Topic title")
		}

		if res.Code != http.StatusOK {
			t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
		}
	})
}

func TestTopicList(t *testing.T) {
	t.Run("renders flash messages from session, if present", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics", nil)
		res := httptest.NewRecorder()

		s.session.SaveFlash("Important flash", req, res)

		s.TopicList(res, req)

		if strings.Contains(res.Body.String(), "<section class=\"flash flash-error\">") == false {
			t.Error("response body should include redirect found link")
		}
	})

	t.Run("renders list of topics successfully", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()

		storage.GetRecentTopicsFunc = func() ([]Topic, error) {
			return []Topic{{Title: "First list title"}, {Title: "Second list title"}}, nil
		}

		s.TopicList(res, req)

		if strings.Contains(res.Body.String(), "First list title") == false {
			t.Error("response body should include first list topic")
		}

		if strings.Contains(res.Body.String(), "Second list title") == false {
			t.Error("response body should include second list topic")
		}
	})
}

func TestMessageCreate(t *testing.T) {
	t.Run("responds with 302 to dashboard if user not logged in", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/3/messages", nil)
		res := httptest.NewRecorder()

		s.MessageCreate(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("responds with 302 to dashboard if saving message fails", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/3/messages", nil)
		res := httptest.NewRecorder()
		s.session.SaveUser(&User{Initials: "AK", Theme: 3}, req, res)
		vars := map[string]string{"id": "3"}
		req = mux.SetURLVars(req, vars)

		storage.CreateMessageFunc = func(m *Message) (*Message, error) {
			return nil, errors.New("something went wrong")
		}

		s.MessageCreate(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("success", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/3/messages", nil)
		res := httptest.NewRecorder()
		s.session.SaveUser(&User{Initials: "AK", Theme: 3}, req, res)
		vars := map[string]string{"id": "3"}
		req = mux.SetURLVars(req, vars)

		s.MessageCreate(res, req)

		assertRedirect("/topics/3", t, res)
	})
}

func TestTopicNew(t *testing.T) {
	t.Run("responds with 302 to dashboard if user not logged in", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/new", nil)
		res := httptest.NewRecorder()

		s.TopicNew(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("renders new topic form for logged in users", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/new", nil)
		res := httptest.NewRecorder()
		s.session.SaveUser(&User{Initials: "AK", Theme: 3}, req, res)

		s.TopicNew(res, req)

		if strings.Contains(res.Body.String(), "<section class=\"new-message-header\">") == false {
			t.Error("response body should include new topic form")
		}
	})
}

func TestTopicCreate(t *testing.T) {
	t.Run("responds with 302 to dashboard if user not logged in", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/new", nil)
		res := httptest.NewRecorder()

		s.TopicCreate(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("responds with 302 to new topic on success", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/new?title=Birdwatchig+tips&content=check+it+out", nil)
		res := httptest.NewRecorder()
		s.session.SaveUser(&User{Initials: "AK", Theme: 3}, req, res)
		topicID := 321

		storage.CreateTopicFunc = func(title string) (*Topic, error) {
			return &Topic{ID: &topicID, Title: title, Messages: &[]Message{}}, nil
		}

		s.TopicCreate(res, req)

		assertRedirect("/topics/321", t, res)
	})
}

func TestJoinShow(t *testing.T) {
	t.Run("responds with 302 to dashboard if user logged in", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/join", nil)
		res := httptest.NewRecorder()
		s.session.SaveUser(&User{Initials: "AK", Theme: 3}, req, res)

		s.JoinShow(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("renders join page", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/join", nil)
		res := httptest.NewRecorder()

		s.JoinShow(res, req)

		if strings.Contains(res.Body.String(), "Join the conversation.") == false {
			t.Error("response body should include join page")
		}
	})
}

func TestJoinCreate(t *testing.T) {
	t.Run("does not save user if user initials invalid", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/join?initials=ABC&theme=1", nil)
		res := httptest.NewRecorder()

		s.JoinCreate(res, req)

		user, _ := s.session.GetUser(req)

		if user != nil {
			t.Error("user should not have been set")
		}

		assertRedirect("/join", t, res)
	})

	t.Run("saves user and redirects home", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/join?initials=AK&theme=3", nil)
		res := httptest.NewRecorder()

		s.JoinCreate(res, req)

		user, _ := s.session.GetUser(req)

		if user.Initials != "AK" {
			t.Error("user should have been set to correct value")
		}

		assertRedirect("/topics", t, res)
	})
}
