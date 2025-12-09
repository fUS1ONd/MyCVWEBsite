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

func TestHandler_togglePostLike(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Like Post Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/posts/1/like", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Like.On("TogglePostLike", mock.Anything, 1, 1).Return(true, nil)

		w := httptest.NewRecorder()
		h.togglePostLike(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unlike Post Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/posts/1/like", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Like.On("TogglePostLike", mock.Anything, 1, 1).Return(false, nil)

		w := httptest.NewRecorder()
		h.togglePostLike(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("No User", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/posts/1/like", nil)
		req = injectParam(req, "id", "1")

		w := httptest.NewRecorder()
		h.togglePostLike(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Post ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/posts/invalid/like", nil)
		req = injectParam(req, "id", "invalid")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.togglePostLike(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/posts/1/like", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Like.On("TogglePostLike", mock.Anything, 1, 1).Return(false, errors.New("db error"))

		w := httptest.NewRecorder()
		h.togglePostLike(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_toggleCommentLike(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Like Comment Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/comments/1/like", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Like.On("ToggleCommentLike", mock.Anything, 1, 1).Return(true, nil)

		w := httptest.NewRecorder()
		h.toggleCommentLike(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("No User", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/comments/1/like", nil)
		req = injectParam(req, "id", "1")

		w := httptest.NewRecorder()
		h.toggleCommentLike(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Comment ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/comments/invalid/like", nil)
		req = injectParam(req, "id", "invalid")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.toggleCommentLike(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_getPostLikesCount(t *testing.T) { //nolint:dupl // similar test pattern to getCommentLikesCount
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts/1/likes", nil)
		req = injectParam(req, "id", "1")

		mocks.Like.On("GetPostLikesCount", mock.Anything, 1).Return(42, nil)

		w := httptest.NewRecorder()
		h.getPostLikesCount(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Post ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts/invalid/likes", nil)
		req = injectParam(req, "id", "invalid")

		w := httptest.NewRecorder()
		h.getPostLikesCount(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts/1/likes", nil)
		req = injectParam(req, "id", "1")

		mocks.Like.On("GetPostLikesCount", mock.Anything, 1).Return(0, errors.New("db error"))

		w := httptest.NewRecorder()
		h.getPostLikesCount(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_getCommentLikesCount(t *testing.T) { //nolint:dupl // similar test pattern to getPostLikesCount
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/comments/1/likes", nil)
		req = injectParam(req, "id", "1")

		mocks.Like.On("GetCommentLikesCount", mock.Anything, 1).Return(15, nil)

		w := httptest.NewRecorder()
		h.getCommentLikesCount(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Comment ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/comments/invalid/likes", nil)
		req = injectParam(req, "id", "invalid")

		w := httptest.NewRecorder()
		h.getCommentLikesCount(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/comments/1/likes", nil)
		req = injectParam(req, "id", "1")

		mocks.Like.On("GetCommentLikesCount", mock.Anything, 1).Return(0, errors.New("db error"))

		w := httptest.NewRecorder()
		h.getCommentLikesCount(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
