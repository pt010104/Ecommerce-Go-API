package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
)

func validateProductInput(ctx context.Context, sc models.Scope, input shop.GetProductFilter) error {
	if input.ShopID == "" && len(input.IDs) == 0 {
		return shop.ErrRequireField
	}
	return nil
}
