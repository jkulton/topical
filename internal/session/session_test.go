package session

import (
	"github.com/jkulton/topical/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser(t *testing.T) {
	t.Run("save and return user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()

		s := NewSession("test")
		u := &models.User{Initials: "JK", Theme: 0}

		s.SaveUser(u, req, res)

		savedUser, _ := s.GetUser(req)

		if savedUser.Initials != "JK" {
			t.Error("user retrieved should match the user saved")
		}
	})

	t.Run("returns error if no user found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)

		s := NewSession("test")
		_, err := s.GetUser(req)

		if err == nil {
			t.Error("expected error to be present")
		}
	})
}

func TestFlashes(t *testing.T) {
	t.Run("save and return flashes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()

		s := NewSession("test")

		s.SaveFlash("1st", req, res)
		s.SaveFlash("2nd", req, res)

		flashes, _ := s.GetFlashes(req, res)

		if flashes[0] != "1st" || flashes[1] != "2nd" {
			t.Error("flashes retrieved should match the flashes saved")
		}
	})

	t.Run("flashes empty after get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/topics/12", nil)
		res := httptest.NewRecorder()

		s := NewSession("test")

		s.SaveFlash("1st", req, res)
		s.SaveFlash("2nd", req, res)

		// First GetFlashes returns flashes
		s.GetFlashes(req, res)
		// Second GetFlashes should return no flashes
		flashes, _ := s.GetFlashes(req, res)

		if len(flashes) != 0 {
			t.Error("no flashes should be returned")
		}
	})
}
