package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal-web-platform/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_listPosts(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts", nil)

		expectedResp := &domain.PostsListResponse{
			Posts:      []domain.Post{{ID: 1, Title: "Test"}},
			TotalCount: 1,
		}
		mocks.Post.On("ListPosts", mock.Anything, mock.AnythingOfType("*domain.ListPostsRequest")).Return(expectedResp, nil)

		w := httptest.NewRecorder()
		h.listPosts(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Internal Error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts", nil)

		mocks.Post.On("ListPosts", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		h.listPosts(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_getPostBySlug(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts/test-slug", nil)
		req = injectParam(req, "slug", "test-slug")

		post := &domain.Post{ID: 1, Slug: "test-slug"}
		mocks.Post.On("GetPostBySlug", mock.Anything, "test-slug").Return(post, nil)

		w := httptest.NewRecorder()
		h.getPostBySlug(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/posts/unknown", nil)
		req = injectParam(req, "slug", "unknown")

		mocks.Post.On("GetPostBySlug", mock.Anything, "unknown").Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		h.getPostBySlug(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_createPost(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		body := `{"title": "New Post", "content": "Content"}`
		req := httptest.NewRequest("POST", "/api/v1/admin/posts", bytes.NewBufferString(body))

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Post.On("CreatePost", mock.Anything, mock.AnythingOfType("*domain.CreatePostRequest"), 1).
			Return(&domain.Post{ID: 1, Title: "New Post"}, nil)

		w := httptest.NewRecorder()
		h.createPost(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestHandler_updatePost(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		body := `{"title": "Updated Post"}`
		req := httptest.NewRequest("PUT", "/api/v1/admin/posts/1", bytes.NewBufferString(body))
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Post.On("UpdatePost", mock.Anything, 1, mock.AnythingOfType("*domain.UpdatePostRequest"), 1, true).
			Return(&domain.Post{ID: 1, Title: "Updated Post"}, nil)

		w := httptest.NewRecorder()
		h.updatePost(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_deletePost(t *testing.T) {
	h, mocks := setupHandler(t)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/admin/posts/1", nil)
		req = injectParam(req, "id", "1")

		user := &domain.User{ID: 1, Role: domain.RoleAdmin}
		ctx := context.WithValue(req.Context(), userContextKey, user)
		req = req.WithContext(ctx)

		mocks.Post.On("DeletePost", mock.Anything, 1, 1, true).Return(nil)

		w := httptest.NewRecorder()
		h.deletePost(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
