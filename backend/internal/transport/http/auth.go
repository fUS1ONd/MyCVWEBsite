package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// authLogin initiates OAuth login flow
func (h *Handler) authLogin(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")

	h.log.Info("auth: initiating login", "provider", providerName)

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		h.log.Error("auth: provider not found", "error", err, "provider", providerName)
		http.Error(w, "provider not found", http.StatusBadRequest)
		return
	}

	// Generate state string
	state := gothic.SetState(r)
	h.log.Debug("auth: state generated", "state", state)

	// Explicitly set state in session to ensure it's saved by StoreInSession later
	sessionStore, _ := gothic.Store.Get(r, gothic.SessionName)
	sessionStore.Values["state"] = state

	// Begin auth flow with provider
	sess, err := provider.BeginAuth(state)
	if err != nil {
		h.log.Error("auth: failed to begin auth", "error", err, "provider", providerName)
		http.Error(w, "failed to begin auth", http.StatusInternalServerError)
		return
	}

	// Explicitly store session in cookie
	if err := gothic.StoreInSession(providerName, sess.Marshal(), r, w); err != nil {
		h.log.Error("auth: failed to store session", "error", err, "provider", providerName)
		http.Error(w, "failed to store session", http.StatusInternalServerError)
		return
	}

	h.log.Debug("auth: session stored in cookie", "provider", providerName)

	// Get auth URL and redirect
	authURL, err := sess.GetAuthURL()
	if err != nil {
		h.log.Error("auth: failed to get auth url", "error", err, "provider", providerName)
		http.Error(w, "failed to get auth url", http.StatusInternalServerError)
		return
	}

	h.log.Info("auth: redirecting to provider", "provider", providerName, "auth_url", authURL)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// authCallback handles OAuth callback
func (h *Handler) authCallback(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")

	h.log.Info("auth: callback received", "provider", providerName, "query_params", r.URL.Query().Encode())

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		h.log.Error("auth: provider not found", "error", err, "provider", providerName)
		http.Error(w, "provider not found", http.StatusBadRequest)
		return
	}

	// Validate state
	urlState := r.URL.Query().Get("state")
	cookieState := gothic.GetState(r)

	h.log.Debug("auth: validating state",
		"url_state", urlState,
		"cookie_state", cookieState,
		"match", urlState == cookieState,
	)

	if urlState != cookieState {
		h.log.Error("auth: state token mismatch", "url_state", urlState, "cookie_state", cookieState, "provider", providerName)
		http.Error(w, "state token mismatch", http.StatusUnauthorized)
		return
	}

	h.log.Info("auth: state validated successfully", "provider", providerName)

	// Retrieve session (includes PKCE verifier)
	sessionStr, err := gothic.GetFromSession(providerName, r)
	if err != nil {
		h.log.Error("auth: failed to get session from store", "error", err, "provider", providerName)
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}

	h.log.Debug("auth: session retrieved from cookie", "provider", providerName, "session_length", len(sessionStr))

	sess, err := provider.UnmarshalSession(sessionStr)
	if err != nil {
		h.log.Error("auth: failed to unmarshal session", "error", err, "provider", providerName)
		http.Error(w, "failed to unmarshal session", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	params := r.URL.Query()

	h.log.Info("auth: starting token exchange", "provider", providerName, "has_code", params.Get("code") != "")

	_, err = sess.Authorize(provider, params)
	if err != nil {
		h.log.Error("auth: token exchange failed", "error", err, "provider", providerName)
		http.Error(w, "authentication failed", http.StatusUnauthorized)
		return
	}

	h.log.Info("auth: token exchange successful", "provider", providerName)

	// Fetch user info
	h.log.Info("auth: fetching user info", "provider", providerName)

	gothUser, err := provider.FetchUser(sess)
	if err != nil {
		h.log.Error("auth: failed to fetch user", "error", err, "provider", providerName)
		http.Error(w, "failed to fetch user data", http.StatusInternalServerError)
		return
	}

	h.log.Info("auth: user info fetched",
		"provider", providerName,
		"user_id", gothUser.UserID,
		"email", gothUser.Email,
	)

	if gothUser.Email == "" {
		h.log.Warn("auth: email is empty", "provider", providerName, "user_id", gothUser.UserID)
	}

	// Login with OAuth
	h.log.Info("auth: creating user session", "provider", providerName)

	user, session, err := h.services.Auth.LoginWithOAuth(r.Context(), gothUser)
	if err != nil {
		h.log.Error("auth: failed to login with oauth", "error", err, "provider", providerName, "email", gothUser.Email)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.log.Info("auth: user session created",
		"provider", providerName,
		"user_id", user.ID,
		"email", user.Email,
		"role", user.Role,
	)

	// Set session cookie
	h.setSessionCookie(w, session.Token)
	h.log.Info("auth: session cookie set", "provider", providerName, "user_id", user.ID)

	// Redirect to frontend blog page
	// User can now access protected routes with the session cookie
	frontendURL := h.cfg.OAuth.FrontendURL + "/blog"

	h.log.Info("auth: authentication complete, redirecting",
		"provider", providerName,
		"user_id", user.ID,
		"redirect_url", frontendURL,
	)

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
