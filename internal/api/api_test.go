package api

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/jkulton/topical/internal/models"
	"github.com/jkulton/topical/internal/session"
	"github.com/jkulton/topical/internal/templates"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockStorage struct {
	GetTopicFunc        func(id int) (*models.Topic, error)
	GetRecentTopicsFunc func() ([]models.Topic, error)
	CreateMessageFunc   func(m *models.Message) (*models.Message, error)
	CreateTopicFunc     func(title string) (*models.Topic, error)
}

func (s *MockStorage) GetTopic(id int) (*models.Topic, error) {
	return s.GetTopicFunc(id)
}

func (s *MockStorage) GetRecentTopics() ([]models.Topic, error) {
	return s.GetRecentTopicsFunc()
}

func (s *MockStorage) CreateMessage(m *models.Message) (*models.Message, error) {
	return s.CreateMessageFunc(m)
}

func (s *MockStorage) CreateTopic(title string) (*models.Topic, error) {
	return s.CreateTopicFunc(title)
}

var (
	testSession   *session.Session
	testTemplates *template.Template
	testStorage   MockStorage
	api           TopicalAPI
)

func setupTests() {
	testSession = session.NewSession("test")
	testTemplates = templates.GenerateTemplates("../../web/views/*.gohtml")
	testStorage = MockStorage{
		GetTopicFunc: func(id int) (*models.Topic, error) {
			return &models.Topic{ID: &id, Title: "First Title"}, nil
		},
		GetRecentTopicsFunc: func() ([]models.Topic, error) {
			return []models.Topic{}, nil
		},
		CreateMessageFunc: func(m *models.Message) (*models.Message, error) {
			return nil, nil
		},
		CreateTopicFunc: func(title string) (*models.Topic, error) {
			return nil, nil
		},
	}

	api = TopicalAPI{testTemplates, &testStorage, testSession}
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

		testStorage.GetTopicFunc = func(id int) (*models.Topic, error) {
			return &models.Topic{}, nil
		}

		api.TopicShow(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("renders topic successfully", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()
		vars := map[string]string{"id": "12"}
		req = mux.SetURLVars(req, vars)

		api.TopicShow(res, req)

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

		api.session.SaveFlash("Important flash", req, res)

		api.TopicList(res, req)

		if strings.Contains(res.Body.String(), "<section class=\"flash flash-error\">") == false {
			t.Error("response body should include redirect found link")
		}
	})

	t.Run("renders list of topics successfully", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()

		testStorage.GetRecentTopicsFunc = func() ([]models.Topic, error) {
			return []models.Topic{{Title: "First list title"}, {Title: "Second list title"}}, nil
		}

		api.TopicList(res, req)

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

		api.MessageCreate(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("responds with 302 to dashboard if saving message fails", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/3/messages", nil)
		res := httptest.NewRecorder()
		api.session.SaveUser(&models.User{Initials: "AK", Theme: 3}, req, res)
		vars := map[string]string{"id": "3"}
		req = mux.SetURLVars(req, vars)

		testStorage.CreateMessageFunc = func(m *models.Message) (*models.Message, error) {
			return nil, errors.New("something went wrong")
		}

		api.MessageCreate(res, req)

		assertRedirect("/topics/3", t, res)
	})

	t.Run("success", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/3/messages", nil)
		res := httptest.NewRecorder()
		api.session.SaveUser(&models.User{Initials: "AK", Theme: 3}, req, res)
		vars := map[string]string{"id": "3"}
		req = mux.SetURLVars(req, vars)

		api.MessageCreate(res, req)

		assertRedirect("/topics/3", t, res)
	})
}

func TestTopicNew(t *testing.T) {
	t.Run("responds with 302 to dashboard if user not logged in", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/new", nil)
		res := httptest.NewRecorder()

		api.TopicNew(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("renders new topic form for logged in users", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/topics/new", nil)
		res := httptest.NewRecorder()
		api.session.SaveUser(&models.User{Initials: "AK", Theme: 3}, req, res)

		api.TopicNew(res, req)

		if strings.Contains(res.Body.String(), "<form class=\"new-message-form") == false {
			t.Error("response body should include new topic form")
		}
	})
}

func TestTopicCreate(t *testing.T) {
	t.Run("responds with 302 to dashboard if user not logged in", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/new", nil)
		res := httptest.NewRecorder()

		api.TopicCreate(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("responds with 302 to new topic on success", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/topics/new?title=Birdwatchig+tips&content=check+it+out", nil)
		res := httptest.NewRecorder()
		api.session.SaveUser(&models.User{Initials: "AK", Theme: 3}, req, res)
		topicID := 321

		testStorage.CreateTopicFunc = func(title string) (*models.Topic, error) {
			return &models.Topic{ID: &topicID, Title: title, Messages: &[]models.Message{}}, nil
		}

		api.TopicCreate(res, req)

		assertRedirect("/topics/321", t, res)
	})
}

func TestJoinShow(t *testing.T) {
	t.Run("responds with 302 to dashboard if user logged in", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/join", nil)
		res := httptest.NewRecorder()
		api.session.SaveUser(&models.User{Initials: "AK", Theme: 3}, req, res)

		api.JoinShow(res, req)

		assertRedirect("/topics", t, res)
	})

	t.Run("renders join page", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodGet, "/join", nil)
		res := httptest.NewRecorder()

		api.JoinShow(res, req)

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

		api.JoinCreate(res, req)

		user, _ := api.session.GetUser(req)

		if user != nil {
			t.Error("user should not have been set")
		}

		assertRedirect("/join", t, res)
	})

	t.Run("saves user and redirects home", func(t *testing.T) {
		setupTests()
		req := httptest.NewRequest(http.MethodPost, "/join?initials=AK&theme=3", nil)
		res := httptest.NewRecorder()

		api.JoinCreate(res, req)

		user, _ := api.session.GetUser(req)

		if user.Initials != "AK" {
			t.Error("user should have been set to correct value")
		}

		assertRedirect("/topics", t, res)
	})
}
