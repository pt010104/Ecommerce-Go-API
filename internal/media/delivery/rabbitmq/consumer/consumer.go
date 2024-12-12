package consumer

import (
	"log"

	rabb "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq"
)

func (c Consumer) Consume() {
	go c.consume(rabb.MediaUploadExchange, rabb.MediaUploadQueueName, c.uploadWorker)
}

func catchPanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic in goroutine: %v", r)
		// Additional error handling or cleanup can go here
	}
}
