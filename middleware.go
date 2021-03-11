package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html"
	"log"
	"net/http"
)

// RequestLoggerMiddleware simply logs simple information during each request
// including the request method and URL path
func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %q", r.Method, html.EscapeString(r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

// ProtectedRouteMiddleware redirects home if a protected route is attempted
// to be accessed without a user present in Context.
func ProtectedRouteMiddleware(protectedRouteNames []string, s *sessions.CookieStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isProtectedRoute := false
			routeName := mux.CurrentRoute(r).GetName()
			_, err := UserFromSession(s, r)

			for _, protected := range protectedRouteNames {
				if routeName == protected {
					isProtectedRoute = true
					break
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
