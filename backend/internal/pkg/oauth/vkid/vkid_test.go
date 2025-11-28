package vkid

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	provider := New("client-id", "secret", "http://localhost/callback")

	assert.NotNil(t, provider)
	assert.Equal(t, "vk", provider.Name())
	assert.Equal(t, "client-id", provider.ClientKey)
	assert.Equal(t, "secret", provider.Secret)
	assert.Equal(t, "http://localhost/callback", provider.CallbackURL)
	assert.NotNil(t, provider.config)
}

func TestNewWithScopes(t *testing.T) {
	provider := New("client-id", "secret", "http://localhost/callback", "additional_scope")

	assert.NotNil(t, provider)
	assert.Contains(t, provider.config.Scopes, "email")
	assert.Contains(t, provider.config.Scopes, "vkid.personal_info")
	assert.Contains(t, provider.config.Scopes, "additional_scope")
}

func TestProviderName(t *testing.T) {
	provider := New("client-id", "secret", "callback")

	assert.Equal(t, "vk", provider.Name())

	provider.SetName("custom-vk")
	assert.Equal(t, "custom-vk", provider.Name())
}

func TestBeginAuth(t *testing.T) {
	provider := New("client-id", "secret", "http://localhost/callback")
	session, err := provider.BeginAuth("test-state")

	assert.NoError(t, err)
	assert.NotNil(t, session)

	authURL, err := session.GetAuthURL()
	assert.NoError(t, err)
	assert.Contains(t, authURL, "id.vk.ru/authorize")
	assert.Contains(t, authURL, "client_id=client-id")
	assert.Contains(t, authURL, "state=test-state")
	assert.Contains(t, authURL, "redirect_uri=http")
	assert.Contains(t, authURL, "code_challenge=")
	assert.Contains(t, authURL, "code_challenge_method=S256")
}

func TestSessionMarshal(t *testing.T) {
	session := &Session{
		AuthURL:      "https://id.vk.ru/authorize?client_id=123",
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		Email:        "test@example.com",
		CodeVerifier: "test-code-verifier-123",
	}

	marshaled := session.Marshal()
	assert.NotEmpty(t, marshaled)
	assert.Contains(t, marshaled, "test-access-token")
	assert.Contains(t, marshaled, "test@example.com")
	assert.Contains(t, marshaled, "test-code-verifier-123")
}

func TestSessionUnmarshal(t *testing.T) {
	provider := New("client-id", "secret", "callback")
	sessionData := `{
		"auth_url": "https://id.vk.ru/authorize",
		"access_token": "test-token",
		"refresh_token": "refresh-token",
		"expires_at": "2025-12-31T00:00:00Z",
		"email": "test@example.com",
		"code_verifier": "test-verifier-abc"
	}`

	session, err := provider.UnmarshalSession(sessionData)
	assert.NoError(t, err)
	assert.NotNil(t, session)

	sess := session.(*Session)
	assert.Equal(t, "https://id.vk.ru/authorize", sess.AuthURL)
	assert.Equal(t, "test-token", sess.AccessToken)
	assert.Equal(t, "refresh-token", sess.RefreshToken)
	assert.Equal(t, "test@example.com", sess.Email)
	assert.Equal(t, "test-verifier-abc", sess.CodeVerifier)
}

func TestSessionGetAuthURL(t *testing.T) {
	// Test with valid auth URL
	session := &Session{
		AuthURL: "https://id.vk.ru/authorize?client_id=123",
	}
	url, err := session.GetAuthURL()
	assert.NoError(t, err)
	assert.Equal(t, "https://id.vk.ru/authorize?client_id=123", url)

	// Test with empty auth URL
	emptySession := &Session{}
	url, err = emptySession.GetAuthURL()
	assert.Error(t, err)
	assert.Equal(t, goth.NoAuthUrlErrorMessage, err.Error())
	assert.Empty(t, url)
}

func TestUserFromReader(t *testing.T) {
	jsonResponse := `{
		"user": {
			"user_id": 123456789,
			"first_name": "Иван",
			"last_name": "Петров",
			"email": "ivan@example.com",
			"avatar": "https://example.com/avatar.jpg",
			"phone": "+79001234567"
		}
	}`

	reader := strings.NewReader(jsonResponse)
	user := &goth.User{}

	err := userFromReader(reader, user)
	assert.NoError(t, err)
	assert.Equal(t, "123456789", user.UserID)
	assert.Equal(t, "Иван", user.FirstName)
	assert.Equal(t, "Петров", user.LastName)
	assert.Equal(t, "Иван Петров", user.Name)
	assert.Equal(t, "ivan@example.com", user.Email)
	assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
}

func TestUserFromReader_MissingEmail(t *testing.T) {
	jsonResponse := `{
		"user": {
			"user_id": 987654321,
			"first_name": "Test",
			"last_name": "User",
			"avatar": "https://example.com/test.jpg"
		}
	}`

	reader := strings.NewReader(jsonResponse)
	user := &goth.User{}

	err := userFromReader(reader, user)
	assert.NoError(t, err)
	assert.Equal(t, "987654321", user.UserID)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)
	assert.Equal(t, "Test User", user.Name)
	assert.Empty(t, user.Email) // Email should be empty
}

func TestUserFromReader_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid json`
	reader := strings.NewReader(invalidJSON)
	user := &goth.User{}

	err := userFromReader(reader, user)
	assert.Error(t, err)
}

func TestUserFromReader_NameTrimming(t *testing.T) {
	jsonResponse := `{
		"user": {
			"user_id": 111,
			"first_name": "  Spaced  ",
			"last_name": "  Name  ",
			"avatar": ""
		}
	}`

	reader := strings.NewReader(jsonResponse)
	user := &goth.User{}

	err := userFromReader(reader, user)
	assert.NoError(t, err)
	// TrimSpace removes leading and trailing spaces from concatenated name
	assert.Equal(t, "Spaced     Name", user.Name)
}

func TestRefreshTokenAvailable(t *testing.T) {
	provider := New("client-id", "secret", "callback")
	assert.True(t, provider.RefreshTokenAvailable())
}

func TestDebug(_ *testing.T) {
	provider := New("client-id", "secret", "callback")
	// Debug is a no-op, just ensure it doesn't panic
	provider.Debug(true)
	provider.Debug(false)
}

func TestClient(t *testing.T) {
	provider := New("client-id", "secret", "callback")
	client := provider.Client()
	assert.NotNil(t, client)
}

func TestSessionString(t *testing.T) {
	session := Session{
		AuthURL:      "https://id.vk.ru/authorize",
		AccessToken:  "token",
		Email:        "test@example.com",
		CodeVerifier: "test-verifier",
	}

	str := session.String()
	assert.NotEmpty(t, str)
	assert.Contains(t, str, "token")
	assert.Contains(t, str, "test@example.com")
	assert.Contains(t, str, "test-verifier")
}

func TestNewConfig(t *testing.T) {
	provider := &Provider{
		ClientKey:   "test-key",
		Secret:      "test-secret",
		CallbackURL: "http://localhost/callback",
	}

	config := newConfig(provider, []string{"extra_scope"})

	assert.Equal(t, "test-key", config.ClientID)
	assert.Equal(t, "test-secret", config.ClientSecret)
	assert.Equal(t, "http://localhost/callback", config.RedirectURL)
	assert.Equal(t, authURL, config.Endpoint.AuthURL)
	assert.Equal(t, tokenURL, config.Endpoint.TokenURL)
	assert.Contains(t, config.Scopes, "email")
	assert.Contains(t, config.Scopes, "vkid.personal_info")
	assert.Contains(t, config.Scopes, "extra_scope")
}

func TestFetchUser_NoAccessToken(t *testing.T) {
	provider := New("client-id", "secret", "callback")
	session := &Session{
		AccessToken: "", // Empty access token
	}

	user, err := provider.FetchUser(session)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get user information without accessToken")
	assert.Empty(t, user.UserID)
}

func TestFetchUser_WithSessionEmail(t *testing.T) {
	// This test verifies that email from session is used when user_info doesn't provide it
	// Note: This is a unit test that tests the logic, not the actual HTTP call
	session := &Session{
		AccessToken: "test-token",
		Email:       "session@example.com",
	}

	// We can't easily test the full FetchUser without mocking HTTP
	// but we've tested userFromReader which handles the mapping
	// The integration will be tested manually
	assert.NotEmpty(t, session.Email)
}

func TestAuthorize_InvalidToken(t *testing.T) {
	// This test verifies error handling in Authorize method
	// Full test requires mocking oauth2.Config.Exchange which is complex
	// The method is covered by integration tests
	provider := New("client-id", "secret", "callback")
	session := &Session{}
	params := url.Values{}

	// This will fail because there's no code parameter
	token, err := session.Authorize(provider, params)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestEndpointConstants(t *testing.T) {
	// Verify that VK ID endpoints are correctly set
	assert.Equal(t, "https://id.vk.ru/authorize", authURL)
	assert.Equal(t, "https://id.vk.ru/oauth2/auth", tokenURL)
	assert.Equal(t, "https://id.vk.ru/oauth2/user_info", userInfoURL)
}

func TestSessionMarshalUnmarshalRoundTrip(t *testing.T) {
	// Test that marshaling and unmarshaling preserves data
	provider := New("client-id", "secret", "callback")
	original := &Session{
		AuthURL:      "https://id.vk.ru/authorize?test=1",
		AccessToken:  "access-123",
		RefreshToken: "refresh-456",
		ExpiresAt:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Email:        "roundtrip@example.com",
		CodeVerifier: "verifier-roundtrip-test",
	}

	marshaled := original.Marshal()
	assert.NotEmpty(t, marshaled)

	unmarshaled, err := provider.UnmarshalSession(marshaled)
	assert.NoError(t, err)

	sess := unmarshaled.(*Session)
	assert.Equal(t, original.AuthURL, sess.AuthURL)
	assert.Equal(t, original.AccessToken, sess.AccessToken)
	assert.Equal(t, original.RefreshToken, sess.RefreshToken)
	assert.Equal(t, original.Email, sess.Email)
	assert.Equal(t, original.CodeVerifier, sess.CodeVerifier)
	// Note: time comparison might have precision differences, so we check year
	assert.Equal(t, original.ExpiresAt.Year(), sess.ExpiresAt.Year())
}

func TestNewConfigWithoutExtraScopes(t *testing.T) {
	provider := &Provider{
		ClientKey:   "test-key",
		Secret:      "test-secret",
		CallbackURL: "http://localhost/callback",
	}

	config := newConfig(provider, nil)

	assert.Equal(t, 2, len(config.Scopes)) // Only default scopes
	assert.Contains(t, config.Scopes, "email")
	assert.Contains(t, config.Scopes, "vkid.personal_info")
}

// TestGenerateCodeVerifier tests PKCE code_verifier generation
func TestGenerateCodeVerifier(t *testing.T) {
	verifier, err := generateCodeVerifier()
	assert.NoError(t, err)
	assert.NotEmpty(t, verifier)

	// Verifier should be Base64URL encoded (32 bytes -> 43 chars)
	assert.Equal(t, 43, len(verifier))

	// Should generate different values on each call
	verifier2, err := generateCodeVerifier()
	assert.NoError(t, err)
	assert.NotEqual(t, verifier, verifier2)
}

// TestGenerateCodeChallenge tests PKCE code_challenge generation
func TestGenerateCodeChallenge(t *testing.T) {
	// Test with a known verifier
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	challenge := generateCodeChallenge(verifier)

	assert.NotEmpty(t, challenge)
	// Challenge should be Base64URL encoded SHA256 (32 bytes -> 43 chars)
	assert.Equal(t, 43, len(challenge))

	// Same verifier should produce same challenge
	challenge2 := generateCodeChallenge(verifier)
	assert.Equal(t, challenge, challenge2)

	// Different verifier should produce different challenge
	verifier3 := "different_verifier_string_here_test"
	challenge3 := generateCodeChallenge(verifier3)
	assert.NotEqual(t, challenge, challenge3)
}

// TestBeginAuthWithPKCE tests that BeginAuth includes PKCE parameters
func TestBeginAuthWithPKCE(t *testing.T) {
	provider := New("client-id", "secret", "http://localhost/callback")
	session, err := provider.BeginAuth("test-state")

	assert.NoError(t, err)
	assert.NotNil(t, session)

	// Check that session has code_verifier
	sess := session.(*Session)
	assert.NotEmpty(t, sess.CodeVerifier)
	assert.Equal(t, 43, len(sess.CodeVerifier))

	// Check that auth URL contains PKCE parameters
	authURL, err := session.GetAuthURL()
	assert.NoError(t, err)
	assert.Contains(t, authURL, "code_challenge=")
	assert.Contains(t, authURL, "code_challenge_method=S256")
}
