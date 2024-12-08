package config

import (
	"github.com/caarlos0/env/v9"
)

type MongoConfig struct {
	Database string `env:"MONGODB_DATABASE"`
	URI      string `env:"MONGODB_URI"`
}

type HTTPServerConfig struct {
	Port int    `env:"PORT" envDefault:"8080"`
	Mode string `env:"API_MODE" envDefault:"debug"`
}

type LoggerConfig struct {
	Level    string `env:"LOG_LEVEL" envDefault:"debug"`
	Mode     string `env:"LOG_MODE" envDefault:"development"`
	Encoding string `env:"LOG_ENCODING" envDefault:"console"`
}

type CloudinaryConfig struct {
	URL string `env:"CLOUDINARY_URL"`
}

type Config struct {
	HTTPServer       HTTPServerConfig
	Logger           LoggerConfig
	Mongo            MongoConfig
	JWT              JWTConfig
	RedisConfig      RedisConfig
	RabbitMQConfig   RabbitMQConfig
	CloudinaryConfig CloudinaryConfig
}
type JWTConfig struct {
	SecretKey string `env:"JWT_SECRET"`
}

type RabbitMQConfig struct {
	URL string `env:"RABBITMQ_URL"`
}

type RedisConfig struct {
	RedisAddr     string `env:"REDIS_HOST"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       string `env:"REDIS_DATABASE"`
	MinIdleConns  int    `env:"REDIS_MIN_IDLE_CONNS"`
	PoolSize      int    `env:"REDIS_POOL_SIZE"`
	PoolTimeout   int    `env:"REDIS_POOL_TIMEOUT"`
	Password      string `env:"REDIS_PASSWORD"`
	DB            int    `env:"REDIS_DATABASE"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
