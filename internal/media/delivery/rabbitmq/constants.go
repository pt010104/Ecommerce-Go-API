package rabbitmq

import "github.com/pt010104/api-golang/pkg/rabbitmq"

const (
	MediaUploadExchangeName = "media_upload_exc"
	MediaUploadQueueName    = "media_upload"
)

var (
	MediaUploadExchange = rabbitmq.ExchangeArgs{
		Name:       MediaUploadExchangeName,
		Type:       rabbitmq.ExchangeTypeFanout,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
	}
)
