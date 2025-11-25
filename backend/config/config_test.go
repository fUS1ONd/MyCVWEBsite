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
`)

	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		tmpfile.Close()
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	os.Setenv("CONFIG_PATH", tmpfile.Name())
	defer os.Unsetenv("CONFIG_PATH")

	cfg := MustLoad()

	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, "localhost:9090", cfg.HTTPServer.Address)
	assert.Equal(t, 1*time.Second, cfg.HTTPServer.Timeout)
	assert.Equal(t, "Test User", cfg.Profile.Name)
}
