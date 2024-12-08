package consumer

import (
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/pkg/log"
	"github.com/pt010104/api-golang/pkg/rabbitmq"
)

// Consumer represents a consumer

type Consumer struct {
	l    log.Logger
	conn *rabbitmq.Connection
	uc   media.UseCase
}

// NewConsumer creates a new consumer
func NewConsumer(l log.Logger, conn *rabbitmq.Connection, uc media.UseCase) Consumer {
	return Consumer{
		l:    l,
		conn: conn,
		uc:   uc,
	}
}
