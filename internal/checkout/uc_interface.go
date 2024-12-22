package checkout

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, sc models.Scope, productIDs []string) (CreateOutput, error)
}
