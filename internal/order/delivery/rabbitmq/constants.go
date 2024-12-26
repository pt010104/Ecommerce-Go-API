package rabbitmq

import "github.com/pt010104/api-golang/pkg/rabbitmq"

const (
	OrderExchangeName = "order_exc"
	OrderQueueName    = "order"
)

var (
	OrderExchange = rabbitmq.ExchangeArgs{
		Name:       OrderExchangeName,
		Type:       rabbitmq.ExchangeTypeFanout,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
)
