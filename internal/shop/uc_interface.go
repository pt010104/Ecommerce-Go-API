package shop

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, sc models.Scope, input CreateShop) (models.Shop, error)
	Get(ctx context.Context, sc models.Scope, input GetShopInput) (GetShopOutput, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	Delete(ctx context.Context, sc models.Scope) error
	Update(ctx context.Context, sc models.Scope, input UpdateInput) (models.Shop, error)

	//Inventory
	CreateInventory(ctx context.Context, sc models.Scope, input CreateInventoryInput) (CreateInventoryOutput, error)
}
