package main

import (
	"fmt"
	"log"
	"github.com/pt010104/api-golang/config"
	"github.com/pt010104/api-golang/internal/appconfig/mongo"
	
)

func main() {
	
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}


	log.Printf("the uri: %s" , cfg.Mongo.URI)
	client, db, err :=mongo.ConnectDB(cfg.Mongo.URI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Disconnect(nil)

	
	log.Printf("Connected to database: %s", db.Name())
	
	fmt.Printf("Mongo URI: %s\n", cfg.Mongo.URI)

}
