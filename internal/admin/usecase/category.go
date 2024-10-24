package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUsecase) CreateCategory(ctx context.Context, sc models.Scope, input admin.CreateCategoryInput) (models.Category, error) {
	if sc.Role == 0 {
		return models.Category{}, admin.ErrNoPermission
	}
	if input.Description == "" || input.Name == "" {
		return models.Category{}, admin.ErrInvalidInput
	}
	cate, err := uc.repo.CreateCategory(ctx, sc, admin.CreateCategoryOption{
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		uc.l.Errorf(ctx, "admin.usecae.CreteCategory.repo.createcategory", err)
		return models.Category{}, err
	}
	return cate, nil
}
