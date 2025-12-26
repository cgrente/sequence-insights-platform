package auth

import (
	"net/http"
	"strings"
)

/*
Package auth contains small, composable HTTP auth middleware.

This project intentionally keeps auth simple:
- If API_KEY is set, clients must provide Authorization: Bearer <API_KEY>
- If API_KEY is empty, auth is disabled (useful for local dev / demos)
*/

// BearerToken returns a middleware that enforces a static bearer token if token != "".
func BearerToken(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.TrimSpace(token) == "" {
				next.ServeHTTP(w, r)
				return
			}
			auth := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if !strings.HasPrefix(auth, prefix) {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			if strings.TrimSpace(strings.TrimPrefix(auth, prefix)) != token {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
