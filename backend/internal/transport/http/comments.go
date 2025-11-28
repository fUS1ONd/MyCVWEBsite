package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"personal-web-platform/internal/domain"

	"github.com/go-chi/chi/v5"
)

// getCommentsByPostSlug handles GET /api/v1/posts/{slug}/comments - get all comments for a post
func (h *Handler) getCommentsByPostSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	comments, err := h.services.Comment.GetCommentsByPostSlug(r.Context(), slug)
	if err != nil {
		h.log.Error("failed to get comments by post slug", "error", err, "slug", slug)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		h.log.Error("failed to encode comments response", "error", err)
	}
}

// createComment handles POST /api/v1/posts/{slug}/comments - create new comment
func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	var req domain.CreateCommentRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode create comment request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context (injected by AuthRequired middleware)
	user := h.getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get post by slug to get post ID
	post, err := h.services.Post.GetPostBySlug(r.Context(), slug)
	if err != nil {
		h.log.Error("failed to get post by slug", "error", err, "slug", slug)
		http.Error(w, "post not found", http.StatusNotFound)
		return
	}

	// Create comment
	comment, err := h.services.Comment.CreateComment(r.Context(), post.ID, &req, user.ID)
	if err != nil {
		h.log.Error("failed to create comment", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		h.log.Error("failed to encode comment response", "error", err)
	}
}

// updateComment handles PUT /api/v1/comments/{id} - update existing comment
func (h *Handler) updateComment(w http.ResponseWriter, r *http.Request) { //nolint:dupl // similar pattern to updatePost
	// Get comment ID from URL
	commentIDStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	var req domain.UpdateCommentRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode update comment request", "error", err)
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

	// Update comment
	comment, err := h.services.Comment.UpdateComment(r.Context(), commentID, &req, user.ID, isAdmin)
	if err != nil {
		h.log.Error("failed to update comment", "error", err, "comment_id", commentID)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		h.log.Error("failed to encode comment response", "error", err)
	}
}

// deleteComment handles DELETE /api/v1/comments/{id} - delete comment
func (h *Handler) deleteComment(w http.ResponseWriter, r *http.Request) {
	// Get comment ID from URL
	commentIDStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
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

	// Delete comment
	if err := h.services.Comment.DeleteComment(r.Context(), commentID, user.ID, isAdmin); err != nil {
		h.log.Error("failed to delete comment", "error", err, "comment_id", commentID)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
