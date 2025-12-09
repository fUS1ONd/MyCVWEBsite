package http

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"personal-web-platform/internal/domain"
	"personal-web-platform/internal/domain/derr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_uploadMedia(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		// Create multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.jpg")
		assert.NoError(t, err)
		_, err = part.Write([]byte("fake image content"))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/admin/media", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		expectedMedia := &domain.MediaFile{
			ID:       1,
			Filename: "test.jpg",
			URL:      "http://example.com/media/1_test.jpg",
		}

		mocks.Media.On("Upload", mock.Anything, mock.Anything, mock.Anything, 1).
			Return(expectedMedia, nil)

		w := httptest.NewRecorder()
		h.uploadMedia(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("No User", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.jpg")
		assert.NoError(t, err)
		_, err = part.Write([]byte("fake image content"))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/admin/media", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		h.uploadMedia(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.jpg")
		assert.NoError(t, err)
		_, err = part.Write([]byte("fake image content"))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/admin/media", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Media.On("Upload", mock.Anything, mock.Anything, mock.Anything, 1).
			Return(nil, derr.ErrValidation)

		w := httptest.NewRecorder()
		h.uploadMedia(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("No File", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/admin/media", strings.NewReader(""))
		req.Header.Set("Content-Type", "multipart/form-data")

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.uploadMedia(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_getMedia(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/media/1", nil)
		req = injectParam(req, "id", "1")

		expectedMedia := &domain.MediaFile{
			ID:       1,
			Filename: "test.jpg",
			URL:      "http://example.com/media/test.jpg",
		}

		mocks.Media.On("GetByID", mock.Anything, 1).Return(expectedMedia, nil)

		w := httptest.NewRecorder()
		h.getMedia(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/media/999", nil)
		req = injectParam(req, "id", "999")

		mocks.Media.On("GetByID", mock.Anything, 999).Return(nil, derr.ErrNotFound)

		w := httptest.NewRecorder()
		h.getMedia(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/media/invalid", nil)
		req = injectParam(req, "id", "invalid")

		w := httptest.NewRecorder()
		h.getMedia(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_deleteMedia(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/media/1", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Media.On("Delete", mock.Anything, 1, 1).Return(nil)

		w := httptest.NewRecorder()
		h.deleteMedia(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/media/1", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 2, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Media.On("Delete", mock.Anything, 1, 2).Return(derr.ErrUnauthorized)

		w := httptest.NewRecorder()
		h.deleteMedia(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/media/999", nil)
		req = injectParam(req, "id", "999")

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Media.On("Delete", mock.Anything, 999, 1).Return(derr.ErrNotFound)

		w := httptest.NewRecorder()
		h.deleteMedia(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("No User", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/media/1", nil)
		req = injectParam(req, "id", "1")

		w := httptest.NewRecorder()
		h.deleteMedia(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandler_listMedia(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/media", nil)

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		expectedFiles := []domain.MediaFile{
			{ID: 1, Filename: "test1.jpg", URL: "http://example.com/media/test1.jpg"},
			{ID: 2, Filename: "test2.jpg", URL: "http://example.com/media/test2.jpg"},
		}

		mocks.Media.On("ListByUploader", mock.Anything, 1).Return(expectedFiles, nil)

		w := httptest.NewRecorder()
		h.listMedia(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("No User", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/media", nil)

		w := httptest.NewRecorder()
		h.listMedia(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/media", nil)

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Media.On("ListByUploader", mock.Anything, 1).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		h.listMedia(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
