package producer

import (
	"context"

	"github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
	"github.com/pt010104/api-golang/pkg/log"
	rmqPkg "github.com/pt010104/api-golang/pkg/rabbitmq"
)

//go:generate mockery --name=Producer
type Producer interface {
	PublishOrderMsg(ctx context.Context, msg rabbitmq.OrderMessage) error

	Run() error

	Close()
}

type implProducer struct {
	l            log.Logger
	conn         rmqPkg.Connection
	orderChannel *rmqPkg.Channel
}

func New(l log.Logger, conn rmqPkg.Connection) Producer {
	return &implProducer{
		l:    l,
		conn: conn,
	}
}
