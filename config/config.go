package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	
)

type MongoConfig struct {
	Database string `env:"MONGODB_DATABASE"` 
	URI      string `env:"MONGODB_URI"`
}
type Config struct {
	Mongo MongoConfig
	
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
