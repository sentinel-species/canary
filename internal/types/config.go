package types

import "time"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logging  LoggingConfig  `yaml:"logging"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	ReadTimeout  int64  `yaml:"read_timeout"`
	WriteTimeout int64  `yaml:"write_timeout"`
	IdleTimeout  int64  `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type DatabaseConfig struct {
	Driver                string        `yaml:"driver"`
	Source                string        `yaml:"source"`
	MaxConnections        int64         `yaml:"max_connections"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_lifetime"`
	ConnectTimeout        time.Duration `yaml:"connect_timeout"`
	ReadTimeout           time.Duration `yaml:"read_timeout"`
	WriteTimeout          time.Duration `yaml:"write_timeout"`
}
