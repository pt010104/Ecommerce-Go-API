package consumer

import (
	"log"

	rabb "github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
)

func (c Consumer) Consume() {
	go c.consume(rabb.OrderExchange, rabb.OrderQueueName, c.orderWorker)
}

func catchPanic() {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic in goroutine: %v", r)
		// Additional error handling or cleanup can go here
	}
}
