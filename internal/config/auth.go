package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/listenarr/listenarr/internal/auth"
	"github.com/spf13/viper"
)

// EnsureAPIKey ensures an API key exists, generating one if needed
func EnsureAPIKey(cfg *Config) error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}
	configFile := filepath.Join(configPath, "config.yml")

	// If API key is already set, validate it
	if cfg.Auth.APIKey != "" {
		if !auth.ValidateAPIKeyFormat(cfg.Auth.APIKey) {
			return fmt.Errorf("invalid API key format in config")
		}
		return nil
	}

	// Generate a new API key
	apiKey, err := auth.GenerateSecureAPIKey()
	if err != nil {
		return fmt.Errorf("failed to generate API key: %w", err)
	}

	// Set the API key in config
	cfg.Auth.APIKey = apiKey
	viper.Set("auth.api_key", apiKey)

	// Save to config file
	if err := viper.WriteConfigAs(configFile); err != nil {
		// If config file doesn't exist, create it
		if os.IsNotExist(err) {
			if err := os.MkdirAll(configPath, 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}
			if err := viper.WriteConfigAs(configFile); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to write config file: %w", err)
		}
	}

	return nil
}

