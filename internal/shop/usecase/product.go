package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUsecase) CreateProduct(ctx context.Context, sc models.Scope, input shop.CreateProductInput) (models.Product, models.Inventory, error) {
	inven, err1 := uc.CreateInventory(ctx, sc, shop.CreateInventoryInput{
		StockLevel:      input.StockLevel,
		ReorderLevel:    input.ReorderLevel,
		ReorderQuantity: input.ReorderQuantity,
	})

	if err1 != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.createproduct.createinventory", err1)
		return models.Product{}, models.Inventory{}, err1
	}

	shopID, err := primitive.ObjectIDFromHex(sc.ShopID)
	if err != nil {
		uc.l.Errorf(ctx, "invalid ShopID format: %v", err)
		return models.Product{}, models.Inventory{}, err
	}
	categoryIDs := make([]primitive.ObjectID, len(input.CategoryID))
	for i, id := range input.CategoryID {
		categoryID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			uc.l.Errorf(ctx, "invalid CategoryID format: %v", err)
			return models.Product{}, models.Inventory{}, err
		}
		categoryIDs[i] = categoryID
	}
	p, err := uc.repo.CreateProduct(ctx, sc, shop.CreateProductOption{
		Name:        input.Name,
		Price:       input.Price,
		InventoryID: inven.ID,
		ShopID:      shopID,
		CategoryID:  categoryIDs,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.product.createproduct: ", err)
		return models.Product{}, models.Inventory{}, err
	}

	return p, inven, nil
}
func (uc implUsecase) DetailProduct(ctx context.Context, sc models.Scope, productID primitive.ObjectID) (shop.DetailProductOutput, error) {
	u, err := uc.repo.Detailproduct(ctx, productID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.detail.product", err)
		return shop.DetailProductOutput{}, err
	}

	inventory, err := uc.repo.DetailInventory(ctx, u.InventoryID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.detail.inven", err)
		return shop.DetailProductOutput{}, err
	}

	categoryIDs := make([]string, len(u.CategoryID))
	for idx, id := range u.CategoryID {
		categoryIDs[idx] = id.Hex()
	}

	shopDetail, err := uc.repo.DetailShop(ctx, sc, u.ShopID.Hex())
	if err != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.detail.shop", err)
		return shop.DetailProductOutput{}, err
	}

	category, err := uc.adminUC.ListCategories(ctx, sc, admins.GetCategoriesFilter{
		IDs: categoryIDs,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.detail.categories", err)
		return shop.DetailProductOutput{}, err
	}

	var categoryNames []string
	for _, cat := range category {
		categoryNames = append(categoryNames, cat.Name)
	}

	output := shop.DetailProductOutput{
		ID:            u.ID.Hex(),
		Name:          u.Name,
		CategoryName:  categoryNames,
		ShopName:      shopDetail.Name,
		InventoryName: inventory.ID.Hex(),
		Price:         u.Price,
	}

	return output, nil
}

func (uc implUsecase) ListProduct(ctx context.Context, sc models.Scope, opt shop.GetProductFilter) ([]models.Product, error) {
	s, err := uc.repo.ListProduct(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.ListProduct: %v", err)
		return []models.Product{}, err
	}

	return s, nil
}
