package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
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
func UserMiddleware(next http.Handler) http.Handler {
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
