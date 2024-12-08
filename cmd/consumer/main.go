package main

import (
	"context"

	"github.com/cloudinary/cloudinary-go"
	"github.com/pt010104/api-golang/config"
	"github.com/pt010104/api-golang/internal/appconfig/mongo"
	"github.com/pt010104/api-golang/internal/appconfig/redis"
	"github.com/pt010104/api-golang/internal/consumer"
	pkgLog "github.com/pt010104/api-golang/pkg/log"
	"github.com/pt010104/api-golang/pkg/rabbitmq"
)

func main() {
	ctx := context.Background()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	l := pkgLog.InitializeZapLogger(pkgLog.ZapConfig{
		Level:    cfg.Logger.Level,
		Mode:     cfg.Logger.Mode,
		Encoding: cfg.Logger.Encoding,
	})

	client, db, err := mongo.ConnectDB(cfg.Mongo.URI)
	if err != nil {
		l.Fatalf(ctx, "Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	conn, err := rabbitmq.Dial(cfg.RabbitMQConfig.URL, true)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	redisClient, err := redis.Connect(cfg.RedisConfig)
	if err != nil {
		panic(err)
	}
	defer redisClient.Disconnect()

	cld, err := cloudinary.NewFromURL(cfg.CloudinaryConfig.URL)
	if err != nil {
		panic(err)
	}

	if err := consumer.NewServer(l, conn, *db, redisClient, *cld).Run(); err != nil {
		l.Fatalf(ctx, "Failed to run consumer server: %v", err)
	}
}
