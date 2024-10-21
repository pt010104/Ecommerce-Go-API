package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/inventory"
	"github.com/pt010104/api-golang/internal/models"
)

func (uc implUseCase) Create(ctx context.Context, sc models.Scope, input inventory.CreateInput) (inventory.CreateOutput, error) {
	//Cho check detail product

	opt := inventory.CreateOption{
		ProductID:       input.ProductID,
		StockLevel:      input.StockLevel,
		ReorderLevel:    input.ReorderLevel,
		ReorderQuantity: input.ReorderQuantity,
	}
	i, err := uc.repo.Create(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "inventory.usecase.implUseCase.Create: %v", err)
		return inventory.CreateOutput{}, err
	}

	return inventory.CreateOutput{
		Inventory: i,
	}, nil

}
