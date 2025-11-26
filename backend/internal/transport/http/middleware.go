package http

import (
	"context"
	"net/http"

	"personal-web-platform/internal/domain"
)

type contextKey string

const userContextKey contextKey = "user"

// AuthRequired middleware checks if user is authenticated
func (h *Handler) AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := h.getSessionToken(r)
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := h.services.Auth.ValidateSession(r.Context(), token)
		if err != nil {
			h.log.Error("failed to validate session", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminRequired middleware checks if user is admin
func (h *Handler) AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r.Context())
		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if user.Role != domain.RoleAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getUserFromContext extracts user from request context
func (h *Handler) getUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}
