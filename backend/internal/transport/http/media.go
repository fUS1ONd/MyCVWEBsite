package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"personal-web-platform/internal/domain/derr"

	"github.com/go-chi/chi/v5"
)

const (
	maxUploadSize = 10 << 20 // 10MB
)

// uploadMedia handles POST /api/v1/admin/media - upload media file
func (h *Handler) uploadMedia(w http.ResponseWriter, r *http.Request) {
	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse multipart form
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		h.log.Error("failed to parse multipart form", "error", err)
		RespondBadRequest(w, "file too large or invalid form data")
		return
	}

	// Get file from request
	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.Error("failed to get file from request", "error", err)
		RespondBadRequest(w, "file is required")
		return
	}
	defer func() {
		_ = file.Close()
	}()

	// Get user from context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		RespondUnauthorized(w, "authentication required")
		return
	}

	// Upload file
	media, err := h.services.Media.Upload(r.Context(), file, header, user.ID)
	if err != nil {
		h.log.Error("failed to upload media", "error", err)
		if errors.Is(err, derr.ErrValidation) {
			RespondBadRequest(w, err.Error())
			return
		}
		RespondInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(media); err != nil {
		h.log.Error("failed to encode media response", "error", err)
	}
}

// getMedia handles GET /api/v1/media/{id} - get media file details
func (h *Handler) getMedia(w http.ResponseWriter, r *http.Request) {
	// Get media ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondBadRequest(w, "invalid media ID")
		return
	}

	// Get media
	media, err := h.services.Media.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get media", "error", err, "id", id)
		if errors.Is(err, derr.ErrNotFound) {
			RespondNotFound(w, "media not found")
			return
		}
		RespondInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(media); err != nil {
		h.log.Error("failed to encode media response", "error", err)
	}
}

// deleteMedia handles DELETE /api/v1/admin/media/{id} - delete media file
func (h *Handler) deleteMedia(w http.ResponseWriter, r *http.Request) {
	// Get media ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondBadRequest(w, "invalid media ID")
		return
	}

	// Get user from context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		RespondUnauthorized(w, "authentication required")
		return
	}

	// Delete media
	if err := h.services.Media.Delete(r.Context(), id, user.ID); err != nil {
		h.log.Error("failed to delete media", "error", err, "id", id)
		if errors.Is(err, derr.ErrNotFound) {
			RespondNotFound(w, "media not found")
			return
		}
		if errors.Is(err, derr.ErrUnauthorized) {
			RespondForbidden(w, "you can only delete your own files")
			return
		}
		RespondInternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listMedia handles GET /api/v1/admin/media - list user's media files
func (h *Handler) listMedia(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := h.getUserFromContext(r.Context())
	if user == nil {
		RespondUnauthorized(w, "authentication required")
		return
	}

	// Get media list
	files, err := h.services.Media.ListByUploader(r.Context(), user.ID)
	if err != nil {
		h.log.Error("failed to list media", "error", err, "userID", user.ID)
		RespondInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		h.log.Error("failed to encode media list response", "error", err)
	}
}
