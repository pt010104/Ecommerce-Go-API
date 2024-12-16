package media

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

//go:generate mockery --name=UseCase
type UseCase interface {
	Upload(ctx context.Context, sc models.Scope, input UploadInput) ([]models.Media, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Media, error)
	List(ctx context.Context, sc models.Scope, input ListInput) ([]models.Media, error)

	ConsumeUploadMsg(ctx context.Context, sc models.Scope, input ConsumeUploadMsgInput) error
}
