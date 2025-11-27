package http

import (
	"net/http"

	"personal-web-platform/internal/domain"
)

// getProfile godoc
//
//	@Summary		Get profile information
//	@Description	Get public CV/profile information (accessible to everyone)
//	@Tags			Profile
//	@Produce		json
//	@Success		200	{object}	Response{data=ProfileResponse}
//	@Failure		500	{object}	ErrorResponse
//	@Router			/api/v1/profile [get]
func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.services.Profile.GetProfile(r.Context())
	if err != nil {
		h.log.Error("failed to get profile", "error", err)
		RespondInternalError(w)
		return
	}

	RespondSuccess(w, ToProfileResponse(profile))
}

// updateProfile godoc
//
//	@Summary		Update profile information
//	@Description	Update CV/profile information (admin only)
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			request	body		UpdateProfileRequest	true	"Profile update data"
//	@Success		200		{object}	Response{data=ProfileResponse}
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		CookieAuth
//	@Router			/api/v1/admin/profile [put]
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
