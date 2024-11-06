package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUsecase) CreateProduct(ctx context.Context, sc models.Scope, input shop.CreateProductInput) (models.Product, models.Inventory, error) {
	util.PrintJson(sc)
	inven, err1 := uc.repo.CreateInventory(ctx, sc, shop.CreateInventoryOption{
		StockLevel:      input.StockLevel,
		ReorderLevel:    input.ReorderLevel,
		ReorderQuantity: input.ReorderQuantity,
	})
	if err1 != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.createproduct.createinventory", err1)
		return models.Product{}, models.Inventory{}, err1
	}

	shop1, err2 := uc.Detail(ctx, sc, "")
	if err2 != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.createproduct.detail shop", err2)
		return models.Product{}, models.Inventory{}, err2
	}

	p, err := uc.repo.CreateProduct(ctx, sc, shop.CreateProductOption{
		Name:        input.Name,
		Price:       input.Price,
		InventoryID: inven.ID,
		ShopID:      shop1.ID,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.product.createproduct: ", err)
		return models.Product{}, models.Inventory{}, err
	}

	return p, inven, nil
}
