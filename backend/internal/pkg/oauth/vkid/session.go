package vkid

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

// Session stores data during the auth process with VK ID
type Session struct {
	AuthURL      string    `json:"auth_url"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Email        string    `json:"email"`
	CodeVerifier string    `json:"code_verifier"` // PKCE code verifier
}

// GetAuthURL returns the URL for the authentication endpoint
func (s *Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize exchanges the authorization code for access token
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)

	// Exchange authorization code for token with PKCE code_verifier
	token, err := p.config.Exchange(
		goth.ContextForClient(p.Client()),
		params.Get("code"),
		oauth2.SetAuthURLParam("code_verifier", s.CodeVerifier),
	)
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("invalid token received from provider")
	}

	// Extract tokens and expiration
	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.ExpiresAt = token.Expiry

	// VK ID may include email in token response extras
	if email, ok := token.Extra("email").(string); ok {
		s.Email = email
	}

	return s.AccessToken, nil
}

// Marshal serializes the session into a string
func (s *Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// String returns a string representation of the session (same as Marshal)
func (s Session) String() string {
	return s.Marshal()
}
