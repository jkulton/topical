package middleware

import (
	"bytes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLogger(t *testing.T) {
	t.Run("logs the method and path of the request", func(t *testing.T) {
		var str bytes.Buffer

		router := mux.NewRouter()
		router.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
		router.Use(RequestLogger)

		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/path", nil)

		log.SetOutput(&str)

		router.ServeHTTP(rw, req)

		if strings.Contains(str.String(), "GET \"/path\"") == false {
			t.Error("logged path didn't match")
		}
	})
}
