package producer

import (
	"context"
	"encoding/json"

	"github.com/pt010104/api-golang/internal/media/delivery/rabbitmq"
	rmqPkg "github.com/pt010104/api-golang/pkg/rabbitmq"
)

func (p *implProducer) PublishUploadMsg(ctx context.Context, msg rabbitmq.UploadMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.mediaUploadChannel.Publish(ctx, rmqPkg.PublishArgs{
		Exchange: rabbitmq.MediaUploadExchangeName,
		Msg: rmqPkg.Publishing{
			Body:        body,
			ContentType: rmqPkg.ContentTypePlainText,
		},
	})
}
