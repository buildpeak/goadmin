package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// ServerConfig is the configuration for the API server.
type ServerConfig struct {
	Port int `json:"port"`
	Auth struct {
		JWTSecret string `json:"jwt_secret"`
	} `json:"auth"`
}

// Config is the configuration for the application.
type Config struct {
	// Env is the environment the application is running in.
	// It could be "development", "staging", "production", etc.
	Env string `json:"env"`

	// DatabaseURL is the URL to the database.
	DatabaseURL string `json:"database_url"`

	// Log is the configuration for the logger.
	Log struct {
		Level  string `json:"level"`
		Pretty bool   `json:"pretty"`
	} `json:"log"`

	// APIServer is the configuration for the API server.
	API ServerConfig `json:"api"`

	Observability ObservabilityConfig `json:"observability"`

	Google struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"google"`
}

// NewConfig returns a new configuration.
func NewConfig() (*Config, error) {
	environ := os.Getenv("ENV")
	if environ == "" {
		environ = "development"
	}

	configDir := "config/api"
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		configDir = dir
	}

	// Load the configuration from the environment.
	knf := koanf.New(".")

	// Load base.toml
	if err := knf.Load(file.Provider(configDir+"/base.toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading base config: %w", err)
	}

	// Load {env}.toml
	if err := knf.Load(file.Provider(configDir+"/"+environ+".toml"), toml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading %s config: %w", environ, err)
	}

	// Load environment variables
	transformer := func(s string) string {
		return strings.ReplaceAll(
			// remove the prefix and turn to lowercase and replace all __ with .
			strings.ToLower(
				strings.Replace(s, "", "", 1), //nolint:gocritic // prefix is empty
			),
			"__",
			".",
		)
	}
	if err := knf.Load(env.Provider("", ".", transformer), nil); err != nil {
		return nil, fmt.Errorf("error loading `env` config: %w", err)
	}

	cfg := Config{}

	// Unmarshal the loaded configuration into the Config struct.
	if err := knf.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
