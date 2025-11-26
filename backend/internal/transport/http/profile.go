package http

import (
	"encoding/json"
	"net/http"

	"personal-web-platform/internal/domain"
)

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.services.Profile.GetProfile(r.Context())
	if err != nil {
		h.log.Error("failed to get profile", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		h.log.Error("failed to encode profile response", "error", err)
	}
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateProfileRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode update profile request", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Update profile
	profile, err := h.services.Profile.UpdateProfile(r.Context(), &req)
	if err != nil {
		h.log.Error("failed to update profile", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		h.log.Error("failed to encode profile response", "error", err)
	}
}
