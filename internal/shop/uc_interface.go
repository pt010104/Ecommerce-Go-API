package shop

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockery --name=UseCase
type UseCase interface {
	Create(ctx context.Context, sc models.Scope, input CreateShop) (models.Shop, error)
	Get(ctx context.Context, sc models.Scope, input GetShopInput) (GetShopOutput, error)
	ListShop(ctx context.Context, sc models.Scope, opt GetShopsFilter) ([]models.Shop, error)
	Detail(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	Delete(ctx context.Context, sc models.Scope) error
	Update(ctx context.Context, sc models.Scope, input UpdateInput) ([]models.Shop, error)

	//Inventory
	CreateInventory(ctx context.Context, sc models.Scope, input CreateInventoryInput) (models.Inventory, error)
	DetailInventory(ctx context.Context, id primitive.ObjectID) (models.Inventory, error)
	ListInventory(ctx context.Context, sc models.Scope, ids []primitive.ObjectID) ([]models.Inventory, error)
	UpdateInventory(ctx context.Context, sc models.Scope, input UpdateInventoryInput) (models.Inventory, error)
	DeleteInventory(ctx context.Context, sc models.Scope, productIDs []primitive.ObjectID) error

	//Product
	CreateProduct(ctx context.Context, sc models.Scope, input CreateProductInput) (models.Product, models.Inventory, error)
	DetailProduct(ctx context.Context, sc models.Scope, productID primitive.ObjectID) (models.Product, models.Inventory, error)
	ListProduct(ctx context.Context, sc models.Scope, opt GetProductFilter) ([]models.Product, error)
}
