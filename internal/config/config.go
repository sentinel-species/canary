package config

import (
	"canary/internal/types"
	"gopkg.in/yaml.v3"
	"os"
)

func DefaultConfig() *types.Config {
	return &types.Config{
		Server: types.ServerConfig{
			Host:         "",
			Port:         "8080",
			ReadTimeout:  10 * 1000000000,
			WriteTimeout: 10 * 1000000000,
			IdleTimeout:  15 * 1000000000,
		},
		Logging: types.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Database: types.DatabaseConfig{
			Driver:                "pgx",
			Source:                "",
			MaxConnections:        1,
			MaxConnectionLifetime: 10 * 1000000000,
			ConnectTimeout:        10 * 1000000000,
			ReadTimeout:           10 * 1000000000,
			WriteTimeout:          10 * 1000000000,
		},
	}
}

var (
	Config = &types.Config{}
)

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, Config)
}
