// Package config provides application configuration structures and loading
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents application configuration
type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer    `yaml:"http_server"`
	Database   Database      `yaml:"database"`
	Profile    ProfileConfig `yaml:"profile"`
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
	return &cfg
}
