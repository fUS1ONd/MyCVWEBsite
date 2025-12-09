package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// togglePostLike handles POST /api/v1/posts/{id}/like - toggle like on a post
func (h *Handler) togglePostLike(w http.ResponseWriter, r *http.Request) { //nolint:dupl // similar pattern to toggleCommentLike
	// Get post ID from URL
	idStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		RespondBadRequest(w, "invalid post ID")
		return
	}

	// Get user from context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		RespondUnauthorized(w, "authentication required")
		return
	}

	// Toggle like
	liked, err := h.services.Like.TogglePostLike(r.Context(), user.ID, postID)
	if err != nil {
		h.log.Error("failed to toggle post like", "error", err, "postID", postID, "userID", user.ID)
		RespondInternalError(w)
		return
	}

	// Get updated likes count
	count, _ := h.services.Like.GetPostLikesCount(r.Context(), postID)

	// Return response with like status
	response := map[string]interface{}{
		"is_liked":    liked,
		"likes_count": count,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode like response", "error", err)
	}
}

// toggleCommentLike handles POST /api/v1/comments/{id}/like - toggle like on a comment
func (h *Handler) toggleCommentLike(w http.ResponseWriter, r *http.Request) { //nolint:dupl // similar pattern to togglePostLike
	// Get comment ID from URL
	idStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		RespondBadRequest(w, "invalid comment ID")
		return
	}

	// Get user from context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		RespondUnauthorized(w, "authentication required")
		return
	}

	// Toggle like
	liked, err := h.services.Like.ToggleCommentLike(r.Context(), user.ID, commentID)
	if err != nil {
		h.log.Error("failed to toggle comment like", "error", err, "commentID", commentID, "userID", user.ID)
		RespondInternalError(w)
		return
	}

	// Get updated likes count
	count, _ := h.services.Like.GetCommentLikesCount(r.Context(), commentID)

	// Return response with like status
	response := map[string]interface{}{
		"is_liked":    liked,
		"likes_count": count,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode like response", "error", err)
	}
}

// getPostLikesCount handles GET /api/v1/posts/{id}/likes - get likes count for a post
func (h *Handler) getPostLikesCount(w http.ResponseWriter, r *http.Request) {
	// Get post ID from URL
	idStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		RespondBadRequest(w, "invalid post ID")
		return
	}

	// Get likes count
	count, err := h.services.Like.GetPostLikesCount(r.Context(), postID)
	if err != nil {
		h.log.Error("failed to get post likes count", "error", err, "postID", postID)
		RespondInternalError(w)
		return
	}

	// Return response
	response := map[string]interface{}{
		"count": count,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode likes count response", "error", err)
	}
}

// getCommentLikesCount handles GET /api/v1/comments/{id}/likes - get likes count for a comment
func (h *Handler) getCommentLikesCount(w http.ResponseWriter, r *http.Request) {
	// Get comment ID from URL
	idStr := chi.URLParam(r, "id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		RespondBadRequest(w, "invalid comment ID")
		return
	}

	// Get likes count
	count, err := h.services.Like.GetCommentLikesCount(r.Context(), commentID)
	if err != nil {
		h.log.Error("failed to get comment likes count", "error", err, "commentID", commentID)
		RespondInternalError(w)
		return
	}

	// Return response
	response := map[string]interface{}{
		"count": count,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode likes count response", "error", err)
	}
}
