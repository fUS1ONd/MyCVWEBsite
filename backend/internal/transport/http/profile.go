package http

import (
	"net/http"

	"personal-web-platform/internal/domain"
)

// getProfile handles profile retrieval
func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.services.Profile.GetProfile(r.Context())
	if err != nil {
		h.log.Error("failed to get profile", "error", err)
		RespondInternalError(w)
		return
	}

	RespondSuccess(w, ToProfileResponse(profile))
}

// updateProfile handles profile updates (admin only)
func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateProfileRequest

	// Decode and validate request
	if !h.DecodeAndValidateRequest(w, r, &req) {
		return
	}

	// Update profile
	profile, err := h.services.Profile.UpdateProfile(r.Context(), &req)
	if err != nil {
		h.log.Error("failed to update profile", "error", err)
		RespondInternalError(w)
		return
	}

	RespondSuccess(w, ToProfileResponse(profile))
}
