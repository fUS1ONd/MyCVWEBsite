package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustLoad(t *testing.T) {
	// Create a temporary config file for testing
	content := []byte(`
env: "test"
http_server:
  address: "localhost:9090"
  timeout: "1s"
  idle_timeout: "10s"
database:
  url: "postgres://user:pass@localhost:5432/db"
profile:
  name: "Test User"
auth:
  session_secret: "test-secret-key"
  session_max_age: "24h"
oauth:
  base_url: "http://localhost:8080"
  google:
    client_id: "test-google-id"
    client_secret: "test-google-secret"
  github:
    client_id: "test-github-id"
    client_secret: "test-github-secret"
  vk:
    client_id: "test-vk-id"
    client_secret: "test-vk-secret"
`)

	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write(content); err != nil {
		if closeErr := tmpfile.Close(); closeErr != nil {
			t.Logf("failed to close temp file: %v", closeErr)
		}
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	require.NoError(t, os.Setenv("CONFIG_PATH", tmpfile.Name()))
	defer func() {
		if err := os.Unsetenv("CONFIG_PATH"); err != nil {
			t.Logf("failed to unset CONFIG_PATH: %v", err)
		}
	}()

	cfg := MustLoad()

	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "localhost:9090", cfg.HTTPServer.Address)
	assert.Equal(t, 1*time.Second, cfg.HTTPServer.Timeout)
	assert.Equal(t, "Test User", cfg.Profile.Name)
}
