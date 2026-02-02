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
		if r.URL.Path == "/health" || r.URL.Path == "/api/v1/auth/signup" || isLoginPath(r) {
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

// isLoginPath checks if the request is a POST to /user/{id}
func isLoginPath(r *http.Request) bool {
	if r.Method != http.MethodPost {
		return false
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	return len(parts) == 2 && parts[0] == "user"
}
