package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUsecase) CreateInventory(ctx context.Context, sc models.Scope, input shop.CreateInventoryInput) (models.Inventory, error) {
	opt := shop.CreateInventoryOption{
		StockLevel:      input.StockLevel,
		ReorderLevel:    input.ReorderLevel,
		ReorderQuantity: input.ReorderQuantity,
	}
	i, err := uc.repo.CreateInventory(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.implUseCase.Create: %v", err)
		return models.Inventory{}, err
	}

	return i, nil
}

func (uc implUsecase) DetailInventory(ctx context.Context, productID primitive.ObjectID) (models.Inventory, error) {
	i, err := uc.repo.DetailInventory(ctx, productID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.implUseCase.Detail: %v", err)
		return models.Inventory{}, err
	}

	return i, nil
}

func (uc implUsecase) ListInventory(ctx context.Context, sc models.Scope, productIDs []primitive.ObjectID) ([]models.Inventory, error) {
	i, err := uc.repo.ListInventory(ctx, sc, productIDs)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.implUseCase.List: %v", err)
		return []models.Inventory{}, err
	}

	return i, nil
}

func (uc implUsecase) UpdateInventory(ctx context.Context, sc models.Scope, input shop.UpdateInventoryInput) (models.Inventory, error) {
	i, err := uc.repo.DetailInventory(ctx, input.ID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.implUseCase.Update: %v", err)
		return models.Inventory{}, err
	}

	opt := shop.UpdateInventoryOption{
		Model:           i,
		StockLevel:      input.StockLevel,
		ReorderLevel:    input.ReorderLevel,
		ReorderQuantity: input.ReorderQuantity,
		ReservedLevel:   input.ReservedLevel,
		SoldQuantity:    input.SoldQuantity,
	}
	ni, err := uc.repo.UpdateInventory(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.implUseCase.Update: %v", err)
		return models.Inventory{}, err
	}

	return ni, nil
}

func (uc implUsecase) DeleteInventory(ctx context.Context, sc models.Scope, productIDs []primitive.ObjectID) error {
	err := uc.repo.DeleteInventory(ctx, sc, productIDs)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.implUseCase.Delete: %v", err)
		return err
	}

	return nil
}
