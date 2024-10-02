package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	JWT   JWTConfig
	Mongo MongoConfig
}

type JWTConfig struct {
	SecretKey string
}

type MongoConfig struct {
	Databse string `env: "MONGO_DATABASE"`
	URI     string `env: "MONGO_URI"`
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	cfg := &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
