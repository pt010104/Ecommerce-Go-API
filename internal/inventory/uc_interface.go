package inventory

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, sc models.Scope, input CreateInput) (CreateOutput, error)
}
