package shop

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type Repo interface {
	CreateShop(ctx context.Context, sc models.Scope, opt CreateShopOption) (models.Shop, error)
	GetShop(ctx context.Context, sc models.Scope, opt GetOption) ([]models.Shop, paginator.Paginator, error)
	DetailShop(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	DeleteShop(ctx context.Context, sc models.Scope) error
	UpdateShop(ctx context.Context, sc models.Scope, option UpdateOption) (models.Shop, error)

	// Inventory
	CreateInventory(ctx context.Context, sc models.Scope, opt CreateInventoryOption) (models.Inventory, error)
}
