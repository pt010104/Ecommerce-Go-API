package consumer

import (
	"context"
	"encoding/json"

	"github.com/pt010104/api-golang/internal/media"
	rabb "github.com/pt010104/api-golang/internal/media/delivery/rabbitmq"
	"github.com/pt010104/api-golang/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c Consumer) uploadWorker(d amqp.Delivery) {
	ctx := context.Background()
	c.l.Info(ctx, "media.delivery.rabbitmq.consumer.uploadWorker.start")

	var msg rabb.UploadMessage
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		c.l.Errorf(ctx, "media.delivery.rabbitmq.consumer.uploadWorker.Unmarshal: %v", err)
		d.Ack(false)
		return
	}

	input := media.ConsumeUploadMsgInput{
		ID:         msg.ID,
		UserID:     msg.UserID,
		ShopID:     msg.ShopID,
		FileName:   msg.FileName,
		File:       msg.File,
		FolderName: msg.FolderName,
	}

	sc := models.Scope{
		UserID: msg.UserID,
		ShopID: msg.ShopID,
	}

	err = c.uc.ConsumeUploadMsg(ctx, sc, input)
	if err != nil {
		c.l.Errorf(ctx, "media.delivery.rabbitmq.consumer.uploadWorker.ConsumeUploadMsg: %v", err)
		d.Ack(false)
		return
	}

	d.Ack(false)
}
