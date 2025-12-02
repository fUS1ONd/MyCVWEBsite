// Package oauth provides OAuth provider initialization
package oauth

import (
	"log/slog"
	"net/http"
	"personal-web-platform/config"
	"personal-web-platform/internal/pkg/oauth/vkid"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

// InitProviders initializes OAuth providers based on config
func InitProviders(cfg *config.Config, log *slog.Logger) {
	store := sessions.NewCookieStore([]byte(cfg.Auth.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   cfg.Auth.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	}
	gothic.Store = store

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

	// VK ID OAuth 2.1
	if cfg.OAuth.VK.Enabled {
		vkProvider := vkid.New(
			cfg.OAuth.VK.ClientID,
			cfg.OAuth.VK.ClientSecret,
			cfg.OAuth.BaseURL+"/auth/vk/callback",
		)
		vkProvider.SetLogger(log)
		providers = append(providers, vkProvider)
	}

	goth.UseProviders(providers...)

	log.Info("oauth: providers initialized",
		"count", len(providers),
		"base_url", cfg.OAuth.BaseURL,
	)

	if cfg.OAuth.VK.Enabled {
		expectedCallback := cfg.OAuth.BaseURL + "/auth/vk/callback"
		log.Info("oauth: VK ID callback URL",
			"url", expectedCallback,
			"note", "Must match VK app registration exactly",
		)

		if strings.Contains(cfg.OAuth.BaseURL, ":8080") {
			log.Warn("oauth: base_url contains port 8080",
				"base_url", cfg.OAuth.BaseURL,
				"note", "VK may not accept callback URLs with explicit ports. Use Nginx proxy instead.",
			)
		}
	}
}
