package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_AuthRequired(t *testing.T) {
	h, mocks := setupHandler(t)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := h.getUserFromContext(r.Context())
		if user != nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError) // Should not happen if middleware works
		}
	})

	t.Run("Authorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid-token"})

		user := &domain.User{ID: 1, Email: "test@example.com"}
		mocks.Auth.On("ValidateSession", mock.Anything, "valid-token").Return(user, nil)

		w := httptest.NewRecorder()
		h.AuthRequired(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("No Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		h.AuthRequired(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Session", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid-token"})

		mocks.Auth.On("ValidateSession", mock.Anything, "invalid-token").Return(nil, nil)

		w := httptest.NewRecorder()
		h.AuthRequired(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: "error-token"})

		mocks.Auth.On("ValidateSession", mock.Anything, "error-token").Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		h.AuthRequired(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_AdminRequired(t *testing.T) {
	h, _ := setupHandler(t)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { //nolint:revive // r unused in simple handler
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Admin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AdminRequired(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Admin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		user := &domain.User{ID: 2, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.AdminRequired(nextHandler).ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
