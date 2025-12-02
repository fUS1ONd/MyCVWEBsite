package vkid

import (
	"encoding/json"
	"errors"
	"fmt"
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
	DeviceID     string    `json:"device_id"`     // VK ID device identifier
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
	code := params.Get("code")

	// Extract device_id from callback parameters
	deviceID := params.Get("device_id")
	if deviceID != "" {
		s.DeviceID = deviceID
		p.log.Info("vkid: device_id received", "device_id_length", len(deviceID))
	} else {
		p.log.Warn("vkid: no device_id in callback parameters")
	}

	p.log.Info("vkid: exchanging code for token", "code_length", len(code))

	// Exchange authorization code for token with PKCE code_verifier and device_id
	token, err := p.config.Exchange(
		goth.ContextForClient(p.Client()),
		code,
		oauth2.SetAuthURLParam("code_verifier", s.CodeVerifier),
		oauth2.SetAuthURLParam("device_id", s.DeviceID),
	)
	if err != nil {
		p.log.Error("vkid: token exchange failed",
			"error", err,
			"error_type", fmt.Sprintf("%T", err),
		)
		return "", err
	}

	if !token.Valid() {
		p.log.Error("vkid: received invalid token")
		return "", errors.New("invalid token received from provider")
	}

	p.log.Info("vkid: token exchange successful",
		"has_access_token", token.AccessToken != "",
		"expires_at", token.Expiry,
	)

	// Extract tokens and expiration
	s.AccessToken = token.AccessToken
	s.RefreshToken = token.RefreshToken
	s.ExpiresAt = token.Expiry

	// VK ID may include email in token response extras
	if email, ok := token.Extra("email").(string); ok {
		s.Email = email
		p.log.Info("vkid: email found in token response", "email", email)
	} else {
		p.log.Debug("vkid: no email in token response extras")
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
