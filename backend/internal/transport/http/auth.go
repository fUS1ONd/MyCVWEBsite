package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

// authLogin initiates OAuth login flow
func (h *Handler) authLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	// Store provider in query for gothic
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(w, r)
}

// authCallback handles OAuth callback
func (h *Handler) authCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	// Store provider in query for gothic
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		h.log.Error("oauth callback failed", "error", err)
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	// Login with OAuth
	_, session, err := h.services.Auth.LoginWithOAuth(r.Context(), gothUser)
	if err != nil {
		h.log.Error("failed to login with oauth", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	h.setSessionCookie(w, session.Token)

	// Redirect to frontend blog page
	// User can now access protected routes with the session cookie
	frontendURL := h.cfg.OAuth.FrontendURL + "/blog"
	http.Redirect(w, r, frontendURL, http.StatusFound)
}

// authMe returns current user info
func (h *Handler) authMe(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.Error("failed to encode user", "error", err)
	}
}

// authLogout logs out the user
func (h *Handler) authLogout(w http.ResponseWriter, r *http.Request) {
	token := h.getSessionToken(r)
	if token != "" {
		if err := h.services.Auth.Logout(r.Context(), token); err != nil {
			h.log.Error("failed to logout", "error", err)
		}
	}

	// Clear session cookie
	h.clearSessionCookie(w)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "logged out"}); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

// setSessionCookie sets session cookie
func (h *Handler) setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.Auth.CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.cfg.Auth.SessionMaxAge.Seconds()),
		HttpOnly: h.cfg.Auth.CookieHTTPOnly,
		Secure:   h.cfg.Auth.CookieSecure,
		SameSite: h.parseSameSite(h.cfg.Auth.CookieSameSite),
	})
}

// clearSessionCookie clears session cookie
func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.cfg.Auth.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// getSessionToken extracts session token from cookie
func (h *Handler) getSessionToken(r *http.Request) string {
	cookie, err := r.Cookie(h.cfg.Auth.CookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// parseSameSite parses SameSite cookie attribute
func (h *Handler) parseSameSite(value string) http.SameSite {
	switch value {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}
