package usecase

import (
	"context"
	"sync"

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
func (uc *implUsecase) DetailProduct(ctx context.Context, sc models.Scope, productID primitive.ObjectID) (shop.DetailProductOutput, error) {
	var (
		u             models.Product
		inventory     models.Inventory
		shopDetail    models.Shop
		category      []models.Category
		categoryIDs   []string
		categoryNames []string
		err           error
		mu            sync.Mutex
		wg            sync.WaitGroup
	)

	u, err = uc.repo.Detailproduct(ctx, productID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.detail.product", err)
		return shop.DetailProductOutput{}, err
	}

	for _, id := range u.CategoryID {
		categoryIDs = append(categoryIDs, id.Hex())
	}

	errCh := make(chan error, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		inv, err := uc.repo.DetailInventory(ctx, u.InventoryID)
		if err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		inventory = inv
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		shop, err := uc.repo.DetailShop(ctx, models.Scope{}, u.ShopID.Hex())
		if err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		shopDetail = shop
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		categories, err := uc.adminUC.ListCategories(ctx, sc, admins.GetCategoriesFilter{IDs: categoryIDs})
		if err != nil {
			errCh <- err
			return
		}
		mu.Lock()
		category = categories
		mu.Unlock()
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			uc.l.Errorf(ctx, "shop.product.usecase.detail", err)
			return shop.DetailProductOutput{}, err
		}
	}

	for _, cat := range category {
		categoryNames = append(categoryNames, cat.Name)
	}

	output := shop.DetailProductOutput{
		ID:            u.ID.Hex(),
		Name:          u.Name,
		CategoryName:  categoryNames,
		Category:      category,
		Inventory:     inventory,
		ShopName:      shopDetail.Name,
		Shop:          shopDetail,
		InventoryName: inventory.ID.Hex(),
		Price:         u.Price,
	}

	return output, nil
}
func (uc implUsecase) ListProduct(ctx context.Context, sc models.Scope, opt shop.GetProductFilter) (shop.ListProductOutput, error) {
	products, err := uc.repo.ListProduct(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.ListProduct: %v", err)
		return shop.ListProductOutput{}, err
	}

	var output shop.ListProductOutput
	var list []shop.DetailProductOutput
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(products))

	for _, p := range products {
		wg.Add(1)

		go func(product models.Product) {
			defer wg.Done()

			detail, err := uc.DetailProduct(ctx, sc, product.ID)
			if err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			list = append(list, detail)
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.ListProduct: %v", err)
			return shop.ListProductOutput{}, err
		}
	}

	output.List = list
	return output, nil
}
func (uc implUsecase) DeleteOneProduct(ctx context.Context, sc models.Scope, ud primitive.ObjectID) error {
	id, err := primitive.ObjectIDFromHex(sc.ShopID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.DeleteProduct.ObjectIDfromhex", err)
		return err
	}
	if id != ud {
		return shop.ErrNoPermissionToDeleteProduct

	}
	p, err1 := uc.repo.Detailproduct(ctx, ud)
	if err1 != nil {
		uc.l.Errorf(ctx, "shop.usecase.DeleteProduct.DetailProduct", err)
		return err
	}
	var wg sync.WaitGroup
	var invenList []primitive.ObjectID
	invenList = append(invenList, p.InventoryID)
	errCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := uc.repo.Delete(ctx, sc, ud)
		if err != nil {
			errCh <- err
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := uc.repo.DeleteInventory(ctx, sc, invenList)
		if err != nil {
			errCh <- err
			return
		}
	}()
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.DeleteProduct: %v", err)
			return err
		}
	}
	return nil

}
func (uc implUsecase) DeleteProduct(ctx context.Context, sc models.Scope, udList []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(udList))

	for _, ud := range udList {
		wg.Add(1)
		go func(ud string) {
			defer wg.Done()

			shpopidstring, err := primitive.ObjectIDFromHex(ud)
			if err != nil {
				uc.l.Errorf(ctx, "shop.usecase.DeleteProduct.ObjectIDFromHex", err)
				errCh <- err
				return
			}
			p, err := uc.repo.Detailproduct(ctx, shpopidstring)
			if err != nil {
				uc.l.Errorf(ctx, "shop.usecase.DeleteProduct.DetailProduct", err)
				errCh <- err
				return
			}
			if sc.ShopID != p.ShopID.Hex() {
				errCh <- shop.ErrNoPermissionToDeleteProduct
				return
			}

			if err := uc.repo.Delete(ctx, sc, shpopidstring); err != nil {
				uc.l.Errorf(ctx, "shop.usecase.DeleteProduct.Delete", err)
				errCh <- err
				return
			}

			invenList := []primitive.ObjectID{p.InventoryID}
			if err := uc.repo.DeleteInventory(ctx, sc, invenList); err != nil {
				uc.l.Errorf(ctx, "shop.usecase.DeleteProduct.DeleteInventory", err)
				errCh <- err
				return
			}
		}(ud)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
