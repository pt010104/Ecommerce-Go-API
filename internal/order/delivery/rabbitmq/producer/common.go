package producer

import (
	"fmt"

	rabb "github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
	rmqPkg "github.com/pt010104/api-golang/pkg/rabbitmq"
)

// Run runs the producer
func (p *implProducer) Run() (err error) {
	if p.orderChannel, err = p.getWriter(rabb.OrderExchange); err != nil {
		fmt.Println("Error when getting writer")
		return
	}
	return
}

// Close closes the producer
func (p *implProducer) Close() {
	p.orderChannel.Close()
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
