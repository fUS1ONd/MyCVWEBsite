package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer    `yaml:"http_server"`
	Database   Database      `yaml:"database"`
	Profile    ProfileConfig `yaml:"profile"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Database struct {
	URL string `yaml:"url" env-required:"true" env:"DATABASE_URL"`
}

type ProfileConfig struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Skills      []string `yaml:"skills"`
	Contacts    Contacts `yaml:"contacts"`
}

type Contacts struct {
	Email    string `yaml:"email"`
	Github   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
}

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
