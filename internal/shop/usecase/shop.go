package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
)

func (uc implUsecase) Register(ctx context.Context, sc models.Scope, input shop.RegisterInput) (models.Shop, error) {
	if input.City == "" || input.Name == "" || input.Street == "" {
		return models.Shop{}, shop.ErrInvalidInput
	}

	return models.Shop{}, nil

}
