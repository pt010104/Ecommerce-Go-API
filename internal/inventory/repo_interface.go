package inventory

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type Repo interface {
	Create(ctx context.Context, sc models.Scope, opt CreateOption) (models.Inventory, error)
}
