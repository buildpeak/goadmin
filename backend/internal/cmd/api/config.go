package api

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	configDir = "config/api"

	ErrLoadBaseConfig = "error loading base config"
	ErrLoadEnvFile    = "error loading %s config"
	ErrLoadEnvVar     = "error loading env var"
)

type Config struct {
	// Env is the environment the application is running in.
	// It could be "development", "staging", "production", etc.
	Env string `koanf:"env"`

	// DatabaseURL is the URL to the database.
	DatabaseURL string `koanf:"database_url"`

	// Log is the configuration for the logger.
	Log struct {
		Level  string `koanf:"level"`
		Pretty bool   `koanf:"pretty"`
	} `koanf:"log"`

	// APIServer is the configuration for the API server.
	APIServer struct {
		Port int `koanf:"port"`
		Auth struct {
			JWTSecret string `koanf:"jwt_secret"`
		} `koanf:"auth"`
	} `koanf:"api_server"`

	Observability struct {
		Collector struct {
			Host               string   `koanf:"host"`
			Port               int      `koanf:"port"`
			Headers            []Header `koanf:"headers"`
			IsInsecure         bool     `koanf:"is_insecure"`
			WithMetricsEnabled bool     `koanf:"with_metrics_enabled"`
		} `koanf:"collector"`
	} `koanf:"observability"`
}

func NewConfig() (*Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Load the configuration from the environment.
	knf := koanf.New(".")

	if err := knf.Load(file.Provider(configDir+"/base.toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading base config: %w", err)
	}

	if err := knf.Load(file.Provider(configDir+"/"+env+".toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading %s config: %w", env)
	}

	return &Config{}, nil
}
