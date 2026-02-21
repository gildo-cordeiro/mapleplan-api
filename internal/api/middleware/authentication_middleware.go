package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtutil "github.com/gildo-cordeiro/mapleplan-api/pkg/jwt"
)

type contextKey string

const ctxUserIDKey contextKey = "userID"

func GetUserIDFromContext(r *http.Request) (string, bool) {
	id, ok := r.Context().Value(ctxUserIDKey).(string)
	return id, ok
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow unauthenticated access to health check and signup/login endpoints
		if r.URL.Path == "/health" || r.URL.Path == "/api/v1/auth/signup" || isPublicPath(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header or cookie
		var token string
		if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			token = strings.TrimPrefix(auth, "Bearer ")
		} else if c, err := r.Cookie("auth_token"); err == nil {
			token = c.Value
		}

		if token == "" {
			http.Error(w, "Unauthorized: no token provided", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := jwtutil.ParseToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), ctxUserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isPublicPath checks if the request path is a public endpoint
func isPublicPath(r *http.Request) bool {
	path := r.URL.Path

	// health (GET /health)
	if r.Method == http.MethodGet && path == "/health" {
		return true
	}

	// auth endpoints (POST /api/v1/auth/signup, POST /api/v1/auth/login)
	if r.Method == http.MethodPost {
		switch path {
		case "/api/v1/auth/signup", "/api/v1/auth/login", "/api/v1/search-partner":
			return true
		}
	}

	return false
}
