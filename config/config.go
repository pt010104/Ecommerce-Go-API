package config

import (
	"github.com/caarlos0/env/v9"
)

type MongoConfig struct {
	Database string `env:"MONGODB_DATABASE"`
	URI      string `env:"MONGODB_URI"`
}

type HTTPServerConfig struct {
	Port int    `env:"APP_PORT" envDefault:"8080"`
	Mode string `env:"API_MODE" envDefault:"debug"`
}

type LoggerConfig struct {
	Level    string `env:"LOG_LEVEL" envDefault:"debug"`
	Mode     string `env:"LOG_MODE" envDefault:"development"`
	Encoding string `env:"LOG_ENCODING" envDefault:"console"`
}

type Config struct {
	HTTPServer HTTPServerConfig
	Logger     LoggerConfig
	Mongo      MongoConfig
	JWT        JWTConfig
}
type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
