package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUsecase) CreateCategory(ctx context.Context, sc models.Scope, input admins.CreateCategoryInput) (models.Category, error) {
	if sc.Role == 0 {
		return models.Category{}, admins.ErrNoPermission
	}
	if input.Description == "" || input.Name == "" {
		return models.Category{}, admins.ErrInvalidInput
	}
	cate, err := uc.repo.CreateCategory(ctx, sc, admins.CreateCategoryOption{
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		uc.l.Errorf(ctx, "admins.usecae.CreteCategory.repo.createcategory", err)
		return models.Category{}, err
	}
	return cate, nil
}

func (uc implUsecase) ListCategories(ctx context.Context, sc models.Scope, filter admins.GetCategoriesFilter) ([]models.Category, error) {
	cates, err := uc.repo.ListCategories(ctx, sc, filter)
	if err != nil {
		uc.l.Errorf(ctx, "admins.usecae.ListCategories.repo.ListCategories", err)
		return nil, err
	}
	return cates, nil
}
