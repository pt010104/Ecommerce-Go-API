package admins

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

type UseCase interface {
	CreateCategory(ctx context.Context, sc models.Scope, input CreateCategoryInput) (models.Category, error)
	ListCategories(ctx context.Context, sc models.Scope, filter GetCategoriesFilter) ([]models.Category, error)

	VerifyShop(ctx context.Context, sc models.Scope, input VerifyShopInput) ([]models.Shop, error)
}
