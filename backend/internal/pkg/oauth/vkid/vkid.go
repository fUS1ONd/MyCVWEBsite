// Package vkid implements VK ID OAuth 2.1 provider for goth
package vkid

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

const (
	authURL = "https://id.vk.ru/authorize"
	//nolint:gosec // This is an OAuth endpoint URL, not a credential
	tokenURL    = "https://id.vk.ru/oauth2/auth"
	userInfoURL = "https://id.vk.ru/oauth2/user_info"
)

// Provider implements goth.Provider for VK ID OAuth 2.1
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
}

// New creates a new VK ID provider
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: "vk",
	}
	p.config = newConfig(p, scopes)
	return p
}

// Name returns the name of the provider
func (p *Provider) Name() string {
	return p.providerName
}

// SetName sets the name of the provider
func (p *Provider) SetName(name string) {
	p.providerName = name
}

// Client returns the HTTP client
func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

// Debug enables or disables debug mode (no-op for this provider)
func (p *Provider) Debug(_ bool) {}

// BeginAuth initiates the OAuth flow with PKCE support
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	// Generate PKCE parameters
	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PKCE verifier: %w", err)
	}

	challenge := generateCodeChallenge(verifier)

	// Build authorization URL with PKCE parameters
	url := p.config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	session := &Session{
		AuthURL:      url,
		CodeVerifier: verifier,
	}
	return session, nil
}

// FetchUser fetches user information from VK ID
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken:  sess.AccessToken,
		Provider:     p.Name(),
		RefreshToken: sess.RefreshToken,
		ExpiresAt:    sess.ExpiresAt,
	}

	if user.AccessToken == "" {
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}

	// Call VK ID user_info endpoint
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return user, err
	}

	// Add Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+sess.AccessToken)

	response, err := p.Client().Do(req)
	if err != nil {
		return user, err
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			p.HTTPClient = nil // Reset client on close error
		}
	}()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("%s responded with a %d trying to fetch user information", p.providerName, response.StatusCode)
	}

	bits, err := io.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	// Parse user info response and store in RawData
	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	// Map VK ID fields to goth.User
	err = userFromReader(bytes.NewReader(bits), &user)
	if err != nil {
		return user, err
	}

	// Use email from session if not in user_info response
	if user.Email == "" && sess.Email != "" {
		user.Email = sess.Email
	}

	return user, nil
}

// RefreshToken refreshes the access token using the refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	ts := p.config.TokenSource(goth.ContextForClient(p.Client()), token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

// RefreshTokenAvailable returns true if refresh token is available
func (p *Provider) RefreshTokenAvailable() bool {
	return true
}

// UnmarshalSession unmarshals a session from a string
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	sess := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(sess)
	return sess, err
}

// newConfig creates a new OAuth2 config for VK ID
func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	// Default scopes for VK ID
	defaultScopes := []string{"email", "vkid.personal_info"}
	c.Scopes = append(c.Scopes, defaultScopes...)

	// Add additional scopes if provided
	if len(scopes) > 0 {
		c.Scopes = append(c.Scopes, scopes...)
	}

	return c
}

// generateCodeVerifier creates a cryptographically secure random string for PKCE
// Returns a Base64URL-encoded string of 32 random bytes (43 characters)
func generateCodeVerifier() (string, error) {
	// Generate 32 random bytes (256 bits of entropy)
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	// Encode to Base64URL without padding
	verifier := base64.RawURLEncoding.EncodeToString(b)
	return verifier, nil
}

// generateCodeChallenge creates SHA256 hash of verifier encoded in Base64URL
// This is used as the code_challenge parameter in the authorization request
func generateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	return challenge
}

// userFromReader parses VK ID user_info response and maps to goth.User
func userFromReader(reader io.Reader, user *goth.User) error {
	// VK ID user_info response structure
	type vkIDUser struct {
		User struct {
			UserID    int64  `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Avatar    string `json:"avatar"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
		} `json:"user"`
	}

	var response vkIDUser
	err := json.NewDecoder(reader).Decode(&response)
	if err != nil {
		return err
	}

	u := response.User

	// Map VK ID fields to goth.User
	user.UserID = strconv.FormatInt(u.UserID, 10)
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.Name = strings.TrimSpace(u.FirstName + " " + u.LastName)
	user.AvatarURL = u.Avatar

	// Email might be in session or user_info response
	if u.Email != "" {
		user.Email = u.Email
	}

	return nil
}
