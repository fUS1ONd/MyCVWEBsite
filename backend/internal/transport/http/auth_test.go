package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_authMe(t *testing.T) {
	h, _ := setupHandler(t)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/me", nil)
		w := httptest.NewRecorder()

		h.authMe(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/me", nil)
		// Inject user into context
		user := &domain.User{ID: 1, Email: "test@example.com", Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.authMe(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp domain.User
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, resp.ID)
		assert.Equal(t, user.Email, resp.Email)
	})
}

func TestHandler_authLogout(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/logout", nil)
		// Add session cookie
		req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid-token"})

		mocks.Auth.On("Logout", mock.Anything, "valid-token").Return(nil)

		w := httptest.NewRecorder()

		h.authLogout(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Check if cookie is cleared
		cookies := w.Result().Cookies()
		assert.NotEmpty(t, cookies)
		assert.Equal(t, "session_id", cookies[0].Name)
		assert.Equal(t, "", cookies[0].Value)
		assert.Equal(t, -1, cookies[0].MaxAge)
	})

	t.Run("No Cookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/logout", nil)
		w := httptest.NewRecorder()

		h.authLogout(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		cookies := w.Result().Cookies()
		assert.NotEmpty(t, cookies)
		assert.Equal(t, -1, cookies[0].MaxAge)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/logout", nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: "error-token"})

		mocks.Auth.On("Logout", mock.Anything, "error-token").Return(errors.New("db error"))

		w := httptest.NewRecorder()

		h.authLogout(w, req)

		assert.Equal(t, http.StatusOK, w.Code) // Logout should still return 200 and clear cookie
	})
}
