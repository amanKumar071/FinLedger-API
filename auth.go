package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"finance_Data/internal/database"
)

// UserContextKey for storing user in request context
type contextKey string

const UserContextKey contextKey = "user"

// User represents authenticated user data
type User struct {
	ID   string
	Role string
}

// AuthMiddleware verifies user authentication only
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("x-user-id")
		if userID == "" {
			http.Error(w, "Unauthorized: x-user-id header missing", http.StatusUnauthorized)
			return
		}

		var role string
		err := database.DB.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
		if err != nil {
			log.Printf("Error fetching user role for ID %s: %v", userID, err)
			http.Error(w, "User not found", http.StatusForbidden)
			return
		}

		// Store authenticated user in context
		user := &User{ID: userID, Role: role}
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns a middleware that checks if user has required role(s)
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserContextKey).(*User)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userRole := strings.TrimSpace(strings.ToLower(user.Role))
			for _, role := range allowedRoles {
				if userRole == strings.ToLower(strings.TrimSpace(role)) {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
		})
	}
}

// Usage in handlers to access user
func GetUserFromContext(r *http.Request) *User {
	user, _ := r.Context().Value(UserContextKey).(*User)
	return user
}
