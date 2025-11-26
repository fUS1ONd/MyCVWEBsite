// Package oauth provides OAuth provider initialization
package oauth

import (
	"personal-web-platform/config"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/vk"
)

// InitProviders initializes OAuth providers based on config
func InitProviders(cfg *config.Config) {
	gothic.Store = sessions.NewCookieStore([]byte(cfg.Auth.SessionSecret))

	providers := make([]goth.Provider, 0, 3)

	// Google OAuth
	if cfg.OAuth.Google.Enabled {
		providers = append(providers, google.New(
			cfg.OAuth.Google.ClientID,
			cfg.OAuth.Google.ClientSecret,
			cfg.OAuth.BaseURL+"/auth/google/callback",
			"email", "profile",
		))
	}

	// GitHub OAuth
	if cfg.OAuth.GitHub.Enabled {
		providers = append(providers, github.New(
			cfg.OAuth.GitHub.ClientID,
			cfg.OAuth.GitHub.ClientSecret,
			cfg.OAuth.BaseURL+"/auth/github/callback",
			"user:email",
		))
	}

	// VK OAuth
	if cfg.OAuth.VK.Enabled {
		providers = append(providers, vk.New(
			cfg.OAuth.VK.ClientID,
			cfg.OAuth.VK.ClientSecret,
			cfg.OAuth.BaseURL+"/auth/vk/callback",
		))
	}

	goth.UseProviders(providers...)
}
