package producer

import (
	"context"

	"github.com/pt010104/api-golang/internal/media/delivery/rabbitmq"
	"github.com/pt010104/api-golang/pkg/log"
	rmqPkg "github.com/pt010104/api-golang/pkg/rabbitmq"
)

//go:generate mockery --name=Producer
type Producer interface {
	PublishUploadMsg(ctx context.Context, msg rabbitmq.UploadMessage) error

	Run() error

	Close()
}

type implProducer struct {
	l                  log.Logger
	conn               rmqPkg.Connection
	mediaUploadChannel *rmqPkg.Channel
}

func New(l log.Logger, conn rmqPkg.Connection) Producer {
	return &implProducer{
		l:    l,
		conn: conn,
	}
}
