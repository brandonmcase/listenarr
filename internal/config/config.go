package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	QBittorrent QBittorrentConfig `mapstructure:"qbittorrent"`
	Jackett   JackettConfig   `mapstructure:"jackett"`
	Plex      PlexConfig      `mapstructure:"plex"`
	Library   LibraryConfig   `mapstructure:"library"`
	Processing ProcessingConfig `mapstructure:"processing"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	APIKey string `mapstructure:"api_key"`
	Enabled bool   `mapstructure:"enabled"`
}

// QBittorrentConfig holds qBittorrent configuration
type QBittorrentConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// JackettConfig holds Jackett configuration
type JackettConfig struct {
	URL   string `mapstructure:"url"`
	APIKey string `mapstructure:"api_key"`
}

// PlexConfig holds Plex configuration
type PlexConfig struct {
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

// LibraryConfig holds library configuration
type LibraryConfig struct {
	Path string `mapstructure:"path"`
}

// ProcessingConfig holds processing configuration
type ProcessingConfig struct {
	TempPath string `mapstructure:"temp_path"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// Set default config path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Read from environment variables
	viper.SetEnvPrefix("LISTENARR")
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, that's okay - use defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return nil, fmt.Errorf("error creating config directory: %w", err)
	}

	// Ensure API key exists if auth is enabled
	if cfg.Auth.Enabled {
		if err := EnsureAPIKey(&cfg); err != nil {
			return nil, fmt.Errorf("error ensuring API key: %w", err)
		}
	}

	return &cfg, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8686)

	// Database defaults
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}
	viper.SetDefault("database.path", filepath.Join(configPath, "listenarr.db"))

	// Auth defaults
	viper.SetDefault("auth.enabled", true)
	// API key will be generated if not set

	// Library defaults
	libraryPath := os.Getenv("LIBRARY_PATH")
	if libraryPath == "" {
		libraryPath = "./library"
	}
	viper.SetDefault("library.path", libraryPath)

	// Processing defaults
	processingPath := os.Getenv("PROCESSING_PATH")
	if processingPath == "" {
		processingPath = "./processing"
	}
	viper.SetDefault("processing.temp_path", processingPath)
}

