package producer

import (
	"fmt"

	rabb "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq"
	rmqPkg "github.com/pt010104/api-golang/pkg/rabbitmq"
)

// Run runs the producer
func (p *implProducer) Run() (err error) {
	if p.mediaUploadChannel, err = p.getWriter(rabb.MediaUploadExchange); err != nil {
		fmt.Println("Error when getting writer")
		return
	}
	return
}

// Close closes the producer
func (p *implProducer) Close() {
	p.mediaUploadChannel.Close()
}

func (p implProducer) getWriter(exchange rmqPkg.ExchangeArgs) (*rmqPkg.Channel, error) {
	ch, err := p.conn.Channel()
	if err != nil {
		fmt.Println("Error when getting channel")
		return nil, err
	}

	err = ch.ExchangeDeclare(exchange)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
