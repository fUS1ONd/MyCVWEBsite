package http

import (
	"encoding/json"
	"net/http"
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
