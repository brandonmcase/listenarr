package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Set a test config path
	testConfigPath := filepath.Join(os.TempDir(), "listenarr-test")
	defer os.RemoveAll(testConfigPath)

	os.Setenv("CONFIG_PATH", testConfigPath)
	defer os.Unsetenv("CONFIG_PATH")

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check defaults
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8686, cfg.Server.Port)
	assert.True(t, cfg.Auth.Enabled)
	assert.NotEmpty(t, cfg.Auth.APIKey)
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	testConfigPath := filepath.Join(os.TempDir(), "listenarr-test-env")
	defer os.RemoveAll(testConfigPath)

	os.Setenv("CONFIG_PATH", testConfigPath)
	os.Setenv("LISTENARR_SERVER_PORT", "9999")
	os.Setenv("LISTENARR_SERVER_HOST", "127.0.0.1")
	defer func() {
		os.Unsetenv("CONFIG_PATH")
		os.Unsetenv("LISTENARR_SERVER_PORT")
		os.Unsetenv("LISTENARR_SERVER_HOST")
	}()

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, 9999, cfg.Server.Port)
}

func TestLoad_LibraryPath(t *testing.T) {
	testConfigPath := filepath.Join(os.TempDir(), "listenarr-test-lib")
	defer os.RemoveAll(testConfigPath)

	testLibPath := filepath.Join(os.TempDir(), "test-library")
	os.Setenv("CONFIG_PATH", testConfigPath)
	os.Setenv("LIBRARY_PATH", testLibPath)
	defer func() {
		os.Unsetenv("CONFIG_PATH")
		os.Unsetenv("LIBRARY_PATH")
		os.RemoveAll(testLibPath)
	}()

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, testLibPath, cfg.Library.Path)
}

func TestEnsureAPIKey_GeneratesNewKey(t *testing.T) {
	testConfigPath := filepath.Join(os.TempDir(), "listenarr-test-keygen")
	defer os.RemoveAll(testConfigPath)

	os.Setenv("CONFIG_PATH", testConfigPath)
	defer os.Unsetenv("CONFIG_PATH")

	cfg := &Config{
		Auth: AuthConfig{
			Enabled: true,
			APIKey:  "",
		},
	}

	err := EnsureAPIKey(cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, cfg.Auth.APIKey)
	assert.True(t, len(cfg.Auth.APIKey) >= 16)
}

func TestEnsureAPIKey_PreservesExistingKey(t *testing.T) {
	testConfigPath := filepath.Join(os.TempDir(), "listenarr-test-preserve")
	defer os.RemoveAll(testConfigPath)

	os.Setenv("CONFIG_PATH", testConfigPath)
	defer os.Unsetenv("CONFIG_PATH")

	existingKey := "existing-api-key-12345"
	cfg := &Config{
		Auth: AuthConfig{
			Enabled: true,
			APIKey:  existingKey,
		},
	}

	err := EnsureAPIKey(cfg)
	require.NoError(t, err)
	assert.Equal(t, existingKey, cfg.Auth.APIKey)
}

