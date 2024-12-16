package media

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, sc models.Scope, opt UploadOption) (models.Media, error)
	Update(ctx context.Context, sc models.Scope, id string, opt UpdateOption) (models.Media, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Media, error)
	List(ctx context.Context, sc models.Scope, opt ListOption) ([]models.Media, error)
}
