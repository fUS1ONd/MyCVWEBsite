package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-web-platform/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_getCommentsByPostSlug(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts/slug/comments", nil)
		req = injectParam(req, "slug", "test-post")

		comments := []domain.Comment{{ID: 1, Content: "Test"}}
		mocks.Comment.On("GetCommentsByPostSlug", mock.Anything, "test-post", 0).Return(comments, nil)

		w := httptest.NewRecorder()
		h.getCommentsByPostSlug(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Slug", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts//comments", nil)
		// empty slug

		w := httptest.NewRecorder()
		h.getCommentsByPostSlug(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_createComment(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		body := `{"content": "New comment"}`
		req := httptest.NewRequest("POST", "/api/v1/posts/test-post/comments", bytes.NewBufferString(body))
		req = injectParam(req, "slug", "test-post")

		user := &domain.User{ID: 1}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		post := &domain.Post{ID: 10, Slug: "test-post"}
		mocks.Post.On("GetPostBySlug", mock.Anything, "test-post", 1).Return(post, nil)
		mocks.Comment.On("CreateComment", mock.Anything, 10, mock.AnythingOfType("*domain.CreateCommentRequest"), 1).
			Return(&domain.Comment{ID: 1, Content: "New comment"}, nil)

		w := httptest.NewRecorder()
		h.createComment(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		body := `{"content": "New comment"}`
		req := httptest.NewRequest("POST", "/api/v1/posts/test-post/comments", bytes.NewBufferString(body))
		req = injectParam(req, "slug", "test-post")

		w := httptest.NewRecorder()
		h.createComment(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandler_updateComment(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		body := `{"content": "Updated content"}`
		req := httptest.NewRequest("PUT", "/api/v1/comments/1", bytes.NewBufferString(body))
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Comment.On("UpdateComment", mock.Anything, 1, mock.AnythingOfType("*domain.UpdateCommentRequest"), 1, false).
			Return(&domain.Comment{ID: 1, Content: "Updated content"}, nil)

		w := httptest.NewRecorder()
		h.updateComment(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Forbidden", func(t *testing.T) {
		body := `{"content": "Updated content"}`
		req := httptest.NewRequest("PUT", "/api/v1/comments/1", bytes.NewBufferString(body))
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 2, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Comment.On("UpdateComment", mock.Anything, 1, mock.Anything, 2, false).
			Return(nil, errors.New("permission denied"))

		w := httptest.NewRecorder()
		h.updateComment(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestHandler_deleteComment(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/comments/1", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleUser}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Comment.On("DeleteComment", mock.Anything, 1, 1, false).Return(nil)

		w := httptest.NewRecorder()
		h.deleteComment(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/comments/abc", nil)
		req = injectParam(req, "id", "abc")

		w := httptest.NewRecorder()
		h.deleteComment(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// injectParam helps inserting route params into chi context
func injectParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	// Chi stores the route context in a context value
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
