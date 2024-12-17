package cart

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Update(ctx context.Context, sc models.Scope, opt UpdateInput) (UpdateOutput, error)
	Add(ctx context.Context, sc models.Scope, input CreateCartInput) error
	GetCart(ctx context.Context, sc models.Scope, opt GetOption) (GetCartOutput, error)
}
