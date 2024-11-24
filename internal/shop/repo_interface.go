package shop

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"

	"github.com/pt010104/api-golang/pkg/paginator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockery --name=Repository
type Repository interface {
	CreateShop(ctx context.Context, sc models.Scope, opt CreateShopOption) (models.Shop, error)
	GetShop(ctx context.Context, sc models.Scope, opt GetOption) ([]models.Shop, paginator.Paginator, error)
	ListShop(ctx context.Context, sc models.Scope, opt GetShopsFilter) ([]models.Shop, error)
	DetailShop(ctx context.Context, sc models.Scope, id string) (models.Shop, error)
	DeleteShop(ctx context.Context, sc models.Scope) error
	UpdateShop(ctx context.Context, sc models.Scope, option UpdateOption) (models.Shop, error)

	// Inventory
	CreateInventory(ctx context.Context, sc models.Scope, opt CreateInventoryOption) (models.Inventory, error)
	DetailInventory(ctx context.Context, ID primitive.ObjectID) (models.Inventory, error)
	ListInventory(ctx context.Context, sc models.Scope, IDs []primitive.ObjectID) ([]models.Inventory, error)
	UpdateInventory(ctx context.Context, sc models.Scope, opt UpdateInventoryOption) (models.Inventory, error)
	DeleteInventory(ctx context.Context, sc models.Scope, IDs []primitive.ObjectID) error

	//product
	CreateProduct(ctx context.Context, sc models.Scope, opt CreateProductOption) (models.Product, error)
	Detailproduct(ctx context.Context, id primitive.ObjectID) (models.Product, error)
	ListProduct(ctx context.Context, sc models.Scope, opt GetProductFilter) ([]models.Product, error)
	Delete(ctx context.Context, sc models.Scope, ud []string) (err error)
	ValidateCategoryIDs(ctx context.Context, categoryIDs []primitive.ObjectID) error
}
