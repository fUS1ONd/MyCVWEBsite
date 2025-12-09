// Package config provides application configuration structures and loading
package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents application configuration
type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer    `yaml:"http_server"`
	Database   Database      `yaml:"database"`
	Profile    ProfileConfig `yaml:"profile"`
	Auth       Auth          `yaml:"auth"`
	OAuth      OAuth         `yaml:"oauth"`
	Media      Media         `yaml:"media"`
	CORS       CORS          `yaml:"cors"`
	RateLimit  RateLimit     `yaml:"rate_limit"`
}

// HTTPServer represents HTTP server configuration
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// Database represents database configuration
type Database struct {
	URL string `yaml:"url" env-required:"true" env:"DATABASE_URL"`
}

// ProfileConfig represents profile configuration loaded from config file
type ProfileConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Skills      []string `yaml:"skills"`
	Contacts    Contacts `yaml:"contacts"`
}

// Contacts represents contact information configuration
type Contacts struct {
	Email    string `yaml:"email"`
	Github   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
}

// Auth represents authentication configuration
type Auth struct {
	SessionSecret  string        `yaml:"session_secret" env:"SESSION_SECRET" env-required:"true"`
	SessionMaxAge  time.Duration `yaml:"session_max_age" env-default:"168h"` // 7 days
	CookieName     string        `yaml:"cookie_name" env-default:"session_id"`
	CookieDomain   string        `yaml:"cookie_domain" env:"COOKIE_DOMAIN" env-default:""`
	CookieSecure   bool          `yaml:"cookie_secure" env-default:"false"`
	CookieHTTPOnly bool          `yaml:"cookie_http_only" env-default:"true"`
	CookieSameSite string        `yaml:"cookie_same_site" env-default:"lax"`
}

// OAuth represents OAuth providers configuration
type OAuth struct {
	BaseURL     string        `yaml:"base_url" env:"OAUTH_BASE_URL" env-required:"true"`
	FrontendURL string        `yaml:"frontend_url" env:"OAUTH_FRONTEND_URL" env-default:"http://localhost:5173"`
	Google      OAuthProvider `yaml:"google"`
	GitHub      OAuthProvider `yaml:"github"`
	VK          OAuthProvider `yaml:"vk"`
}

// OAuthProvider represents individual OAuth provider configuration
type OAuthProvider struct {
	ClientID     string `yaml:"client_id" env-required:"true"`
	ClientSecret string `yaml:"client_secret" env-required:"true"`
	Enabled      bool   `yaml:"enabled" env-default:"false"`
}

// Media represents media file handling configuration
type Media struct {
	UploadPath string `yaml:"upload_path" env:"MEDIA_UPLOAD_PATH" env-default:"./uploads"`
	BaseURL    string `yaml:"base_url" env:"MEDIA_BASE_URL" env-default:"http://localhost:8080"`
}

// CORS represents CORS configuration
type CORS struct {
	Enabled        bool     `yaml:"enabled" env-default:"true"`
	AllowedOrigins []string `yaml:"allowed_origins" env-default:"*"`
	AllowedMethods []string `yaml:"allowed_methods" env-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders []string `yaml:"allowed_headers" env-default:"Accept,Authorization,Content-Type,X-CSRF-Token"`
	ExposedHeaders []string `yaml:"exposed_headers" env-default:"Link"`
	AllowedCreds   bool     `yaml:"allowed_credentials" env-default:"true"`
	MaxAge         int      `yaml:"max_age" env-default:"300"`
}

// RateLimit represents rate limiting configuration
type RateLimit struct {
	Enabled       bool `yaml:"enabled" env-default:"true"`
	RequestsLimit int  `yaml:"requests_limit" env-default:"100"` // requests per window
	WindowSeconds int  `yaml:"window_seconds" env-default:"60"`  // time window in seconds
}

// MustLoad loads configuration from file or panics if unable to load
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// Apply defaults and validate configuration
	if cfg.OAuth.FrontendURL == "" {
		cfg.OAuth.FrontendURL = cfg.OAuth.BaseURL
	}

	// Validate OAuth URLs
	if !strings.HasPrefix(cfg.OAuth.BaseURL, "http://") &&
		!strings.HasPrefix(cfg.OAuth.BaseURL, "https://") {
		log.Fatalf("oauth.base_url must start with http:// or https://, got: %s", cfg.OAuth.BaseURL)
	}

	if !strings.HasPrefix(cfg.OAuth.FrontendURL, "http://") &&
		!strings.HasPrefix(cfg.OAuth.FrontendURL, "https://") {
		log.Fatalf("oauth.frontend_url must start with http:// or https://, got: %s", cfg.OAuth.FrontendURL)
	}

	return &cfg
}
