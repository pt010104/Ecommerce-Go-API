package main

import (
	"github.com/pt010104/api-golang/config"
	"github.com/pt010104/api-golang/internal/appconfig/mongo"
	"github.com/pt010104/api-golang/internal/appconfig/redis"
	httpServer "github.com/pt010104/api-golang/internal/httpserver"
	pkgLog "github.com/pt010104/api-golang/pkg/log"

	"log"
)

func main() {
	//t
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, db, err := mongo.ConnectDB(cfg.Mongo.URI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Disconnect(nil)

	log.Printf("Connected to database: %s", db.Name())

	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	redisClient, err := redis.Connect(cfg.RedisConfig)
	if err != nil {
		panic(err)
	}
	defer redisClient.Disconnect()

	log.Println("Sending verification email...")
	srv := httpServer.New(l, httpServer.Config{
		Port:         cfg.HTTPServer.Port,
		JWTSecretKey: cfg.JWT.SecretKey,
		Mode:         cfg.HTTPServer.Mode,
		Database:     *db,
		Redis:        redisClient,
	})

	if err := srv.Run(); err != nil {

		panic(err)
	}

}
