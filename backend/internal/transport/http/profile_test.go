package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_getProfile(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		h, mocks := setupHandler(t)
		req := httptest.NewRequest("GET", "/api/v1/profile", nil)

		profile := &domain.Profile{Name: "Test Name"}
		mocks.Profile.On("GetProfile", mock.Anything).Return(profile, nil)

		w := httptest.NewRecorder()
		h.getProfile(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error", func(t *testing.T) {
		h, mocks := setupHandler(t)
		req := httptest.NewRequest("GET", "/api/v1/profile", nil)

		mocks.Profile.On("GetProfile", mock.Anything).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		h.getProfile(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_updateProfile(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		h, mocks := setupHandler(t)
		body := `{"name": "New Name", "description": "Desc", "activity": "Work", "contacts": {"email": "test@test.com"}}`
		req := httptest.NewRequest("PUT", "/api/v1/admin/profile", bytes.NewBufferString(body))

		// Note: AdminRequired middleware check is done before this handler
		// Here we assume it passed and check handler logic

		mocks.Profile.On("UpdateProfile", mock.Anything, mock.AnythingOfType("*domain.UpdateProfileRequest")).
			Return(&domain.Profile{Name: "New Name"}, nil)

		w := httptest.NewRecorder()
		h.updateProfile(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
