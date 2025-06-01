package types

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
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
