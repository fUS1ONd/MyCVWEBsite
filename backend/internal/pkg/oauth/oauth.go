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
	enabledProviders := []string{}

	// Google OAuth
	if cfg.OAuth.Google.Enabled {
		providers = append(providers, google.New(
			cfg.OAuth.Google.ClientID,
			cfg.OAuth.Google.ClientSecret,
			cfg.OAuth.BaseURL+"/auth/google/callback",
			"email", "profile",
		))
		enabledProviders = append(enabledProviders, "google")
		log.Info("oauth: Google provider enabled",
			"callback_url", cfg.OAuth.BaseURL+"/auth/google/callback",
		)
	}

	// GitHub OAuth
	if cfg.OAuth.GitHub.Enabled {
		providers = append(providers, github.New(
			cfg.OAuth.GitHub.ClientID,
			cfg.OAuth.GitHub.ClientSecret,
			cfg.OAuth.BaseURL+"/auth/github/callback",
			"user:email",
		))
		enabledProviders = append(enabledProviders, "github")
		log.Info("oauth: GitHub provider enabled",
			"callback_url", cfg.OAuth.BaseURL+"/auth/github/callback",
		)
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
		enabledProviders = append(enabledProviders, "vk")
		log.Info("oauth: VK ID provider enabled",
			"callback_url", cfg.OAuth.BaseURL+"/auth/vk/callback",
		)

		if strings.Contains(cfg.OAuth.BaseURL, ":8080") {
			log.Warn("oauth: base_url contains port 8080",
				"base_url", cfg.OAuth.BaseURL,
				"note", "VK may not accept callback URLs with explicit ports. Use Nginx proxy instead.",
			)
		}
	}

	if len(providers) == 0 {
		log.Warn("oauth: no providers enabled")
		return
	}

	goth.UseProviders(providers...)

	log.Info("oauth: initialization complete",
		"count", len(providers),
		"providers", strings.Join(enabledProviders, ", "),
		"base_url", cfg.OAuth.BaseURL,
		"frontend_url", cfg.OAuth.FrontendURL,
	)
}
