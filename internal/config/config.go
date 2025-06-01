package config

import "canary/internal/types"

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
	}
}
