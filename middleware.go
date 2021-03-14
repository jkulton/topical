package main

import (
	"html"
	"log"
	"net/http"
)

// RequestLoggerMiddleware logs the method and URL path for each request
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %q", r.Method, html.EscapeString(r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
