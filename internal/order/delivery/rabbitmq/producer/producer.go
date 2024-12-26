package producer

import (
	"context"
	"encoding/json"

	"github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
	rmqPkg "github.com/pt010104/api-golang/pkg/rabbitmq"
)

func (p *implProducer) PublishOrderMsg(ctx context.Context, msg rabbitmq.OrderMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.orderChannel.Publish(ctx, rmqPkg.PublishArgs{
		Exchange: rabbitmq.OrderExchangeName,
		Msg: rmqPkg.Publishing{
			Body:        body,
			ContentType: rmqPkg.ContentTypePlainText,
		},
	})
}
