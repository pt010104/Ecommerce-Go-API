package mongo

import (
	"context"
	"fmt"
	"github.com/pt010104/api-golang/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectDB(uri string) (*mongo.Client, *mongo.Database, error) {

	clientOptions := options.Client().ApplyURI(uri)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully")

	db := client.Database(cfg.Mongo.Database)

	return client, db, nil
}
