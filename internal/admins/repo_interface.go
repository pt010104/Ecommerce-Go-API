package admins

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

//go:generate mockery --name=Repository
type Repository interface {
	CreateCategory(ctx context.Context, sc models.Scope, opt CreateCategoryOption) (models.Category, error)
	ListCategories(ctx context.Context, sc models.Scope, filter GetCategoriesFilter) ([]models.Category, error)
}
