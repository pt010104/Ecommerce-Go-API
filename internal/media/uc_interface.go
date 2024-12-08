package media

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Upload(ctx context.Context, sc models.Scope, input UploadInput) error
	Detail(ctx context.Context, sc models.Scope, id string) (models.Media, error)

	ConsumeUploadMsg(ctx context.Context, sc models.Scope, input ConsumeUploadMsgInput) error
}
