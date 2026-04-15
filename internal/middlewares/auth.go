package middlewares

import (
	"context"
	"expense-tracker/internal/service"
	"net/http"
)

// Context key type to avoid coalition
type contextKey string

const UserIDKey contextKey = "user_id"

// AuthMiddleware validates the access token and injects user_id into context
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			// inject user_id into context
			ctx := context.WithValue(r.Context(), UserIDKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
		)
	}
}

// helpers for handlers to extract user_id
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
