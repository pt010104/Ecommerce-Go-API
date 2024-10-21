package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
)

func (uc implUsecase) CreateInventory(ctx context.Context, sc models.Scope, input shop.CreateInventoryInput) (shop.CreateInventoryOutput, error) {

	opt := shop.CreateInventoryOption{
		ProductID:       input.ProductID,
		StockLevel:      input.StockLevel,
		ReorderLevel:    input.ReorderLevel,
		ReorderQuantity: input.ReorderQuantity,
	}
	i, err := uc.repo.CreateInventory(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.implUseCase.Create: %v", err)
		return shop.CreateInventoryOutput{}, err
	}

	return shop.CreateInventoryOutput{
		Inventory: i,
	}, nil

}
