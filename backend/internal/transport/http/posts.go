package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"personal-web-platform/internal/domain"

	"github.com/go-chi/chi/v5"
)

// listPosts handles GET /api/v1/posts - list all posts with pagination and filters
func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := &domain.ListPostsRequest{}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil {
			req.Page = page
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			req.Limit = limit
		}
	}

	if publishedStr := r.URL.Query().Get("published"); publishedStr != "" {
		published := publishedStr == "true"
		req.Published = &published
	}

	if user := h.getUserFromContext(r.Context()); user != nil {
		req.UserID = user.ID
	}

	// Get posts
	response, err := h.services.Post.ListPosts(r.Context(), req)
	if err != nil {
		h.log.Error("failed to list posts", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode posts response", "error", err)
	}
}

// getPostBySlug handles GET /api/v1/posts/{slug} - get post details by slug
func (h *Handler) getPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	var userID int
	if user := h.getUserFromContext(r.Context()); user != nil {
		userID = user.ID
	}

	post, err := h.services.Post.GetPostBySlug(r.Context(), slug, userID)
	if err != nil {
		h.log.Error("failed to get post by slug", "error", err, "slug", slug)
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		h.log.Error("failed to encode post response", "error", err)
	}
}

// createPost handles POST /api/v1/admin/posts - create new post
func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	var req domain.CreatePostRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode create post request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context (injected by AuthRequired middleware)
	user := h.getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Create post
	post, err := h.services.Post.CreatePost(r.Context(), &req, user.ID)
	if err != nil {
		h.log.Error("failed to create post", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		h.log.Error("failed to encode post response", "error", err)
	}
}

// updatePost handles PUT /api/v1/admin/posts/{id} - update existing post
func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) { //nolint:dupl // similar pattern to updateComment
	// Get post ID from URL
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	var req domain.UpdatePostRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode update post request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	isAdmin := user.Role == domain.RoleAdmin

	// Update post
	post, err := h.services.Post.UpdatePost(r.Context(), postID, &req, user.ID, isAdmin)
	if err != nil {
		h.log.Error("failed to update post", "error", err, "post_id", postID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		h.log.Error("failed to encode post response", "error", err)
	}
}

// deletePost handles DELETE /api/v1/admin/posts/{id} - delete post
func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	// Get post ID from URL
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	// Get user from context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	isAdmin := user.Role == domain.RoleAdmin

	// Delete post
	if err := h.services.Post.DeletePost(r.Context(), postID, user.ID, isAdmin); err != nil {
		h.log.Error("failed to delete post", "error", err, "post_id", postID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
