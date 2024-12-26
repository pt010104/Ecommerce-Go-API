package consumer

import (
	"context"
	"encoding/json"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	rabb "github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c Consumer) orderWorker(d amqp.Delivery) {
	ctx := context.Background()
	c.l.Info(ctx, "order.delivery.rabbitmq.consumer.orderWorker.start")

	var msg rabb.OrderMessage
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		c.l.Errorf(ctx, "order.delivery.rabbitmq.consumer.orderWorker.Unmarshal: %v", err)
		d.Ack(false)
		return
	}

	input := order.ConsumeOrderMsgInput{
		OrderID:    msg.OrderID,
		CheckoutID: msg.CheckoutID,
		UserID:     msg.UserID,
	}

	sc := models.Scope{
		UserID: msg.UserID,
	}

	err = c.uc.ConsumeOrderMsg(ctx, sc, input)
	if err != nil {
		c.l.Errorf(ctx, "order.delivery.rabbitmq.consumer.orderWorker.ConsumeOrderMsg: %v", err)
		d.Ack(false)
		return
	}

	d.Ack(false)
}
