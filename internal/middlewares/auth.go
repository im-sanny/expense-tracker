package middlewares

import (
	"context"
	"expense-tracker/internal/service"
	"net/http"
)

// Context key type to avoid coalition
type contextKey string

const UserContextKey contextKey = "user_id"

// AuthMiddleware validates the access token and injects user_id into context
func AuthMiddleware(authService *service.Auth, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Missing access token", http.StatusUnauthorized)
			return
		}

		userId, err := authService.ValidateAccessToken(cookie.Value)
		if err != nil {
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}

		// Add user_id to request context
		ctx := context.WithValue(r.Context(), UserContextKey, userId)
		next(w, r.WithContext(ctx))
	}
}

// GenerateIDFromContext extracts user_id from context in your handlers
func GenerateIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserContextKey).(string)
	return userID, ok
}
