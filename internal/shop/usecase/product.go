package usecase

import (
	"context"
	"fmt"
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
	if input.MediaID != "" {

		_, err := uc.mediaUC.Detail(ctx, sc, input.MediaID)
		if err != nil {
			uc.l.Errorf(ctx, "shop.product.usecase.createproduct.detailmedia", err)
			return models.Product{}, models.Inventory{}, err
		}
	}
	media_id, err := primitive.ObjectIDFromHex(input.MediaID)
	if err != nil {
		uc.l.Errorf(ctx, "invalid MediaID format: %v", err)
		return models.Product{}, models.Inventory{}, err
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
		MediaID:     media_id,
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
		avatar        *models.Media
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

	wg.Add(4)

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
		if u.MediaID != primitive.NilObjectID {
			avatar1, err := uc.mediaUC.Detail(ctx, sc, u.MediaID.Hex())
			if err != nil {
				errCh <- err
				return
			}
			avatar = &avatar1
		}

		mu.Lock()

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
	var mediaID, url string

	if avatar != nil {
		mediaID = avatar.ID.Hex()
		url = avatar.URL
	}
	output := shop.DetailProductOutput{
		ID:           u.ID.Hex(),
		Name:         u.Name,
		CategoryName: categoryNames,
		Category:     category,
		Inventory:    inventory,
		MediaID:      mediaID,
		URL:          url,
		Shop:         shopDetail,

		Price: u.Price,
	}

	return output, nil
}

func (uc implUsecase) ListProduct(ctx context.Context, sc models.Scope, opt shop.GetProductFilter) (shop.ListProductOutput, error) {
	var (
		products []models.Product
		s        models.Shop
	)

	var wg sync.WaitGroup
	var wgErr error
	wg.Add(2)

	go func() {
		defer wg.Done()
		var err error
		products, err = uc.repo.ListProduct(ctx, sc, opt)
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.ListProduct.repo.ListProduct: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		s, err = uc.repo.DetailShop(ctx, models.Scope{}, opt.ShopID)
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.ListProduct.repo.DetailShop: %v", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()

	if wgErr != nil {
		return shop.ListProductOutput{}, wgErr
	}

	categoryIDSet := make(map[string]struct{})
	for _, p := range products {
		for _, catID := range p.CategoryID {
			categoryIDSet[catID.Hex()] = struct{}{}
		}
	}

	var categoryIDs []string
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	categories, err := uc.adminUC.ListCategories(ctx, models.Scope{}, admins.GetCategoriesFilter{
		IDs: categoryIDs,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.ListProduct.adminUC.ListCategories: %v", err)
		return shop.ListProductOutput{}, err
	}

	categoryMap := make(map[string]models.Category)
	for _, cate := range categories {
		categoryMap[cate.ID.Hex()] = cate
	}

	var productsOutPut []shop.ProductOutPutItem
	for _, p := range products {
		var cates []models.Category
		for _, catID := range p.CategoryID {
			if cate, ok := categoryMap[catID.Hex()]; ok {
				cates = append(cates, cate)
			}
		}

		item := shop.ProductOutPutItem{
			P:     p,
			Inven: p.InventoryID.Hex(),
			Cate:  cates,
		}

		productsOutPut = append(productsOutPut, item)
	}

	return shop.ListProductOutput{
		Products: productsOutPut,
		Shop:     s,
	}, nil
}

func (uc implUsecase) DeleteProduct(ctx context.Context, sc models.Scope, udList []string) error {

	err := uc.repo.Delete(ctx, sc, udList)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase. DeleteProduct.repoDelete: %v", err)
		return err
	}
	return nil
}
func (uc implUsecase) GetProduct(ctx context.Context, sc models.Scope, input shop.GetProductOption) (shop.GetProductOutput, error) {
	opt := shop.GetProductOption{
		GetProductFilter: input.GetProductFilter,
		PagQuery:         input.PagQuery,
	}

	s, pag, err := uc.repo.GetProduct(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.GetProduct: %v", err)
		return shop.GetProductOutput{}, err
	}
	//print opt.IDs
	fmt.Println("opt.IDs", opt.IDs)
	categoryIDSet := make(map[string]struct{})
	for _, p := range s {
		for _, catID := range p.CategoryID {
			categoryIDSet[catID.Hex()] = struct{}{}
		}
	}

	var categoryIDs []string
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	categories, err := uc.adminUC.ListCategories(ctx, models.Scope{}, admins.GetCategoriesFilter{
		IDs: categoryIDs,
	})

	categoryMap := make(map[string]models.Category)
	for _, cate := range categories {
		categoryMap[cate.ID.Hex()] = cate
	}
	var list []shop.ProductOutPutItem
	for _, p := range s {
		var cates []models.Category
		for _, catID := range p.CategoryID {
			if cate, ok := categoryMap[catID.Hex()]; ok {
				cates = append(cates, cate)
			}
		}
		avatar, err := uc.mediaUC.Detail(ctx, models.Scope{}, p.MediaID.Hex())

		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetProduct: %v", err)
			return shop.GetProductOutput{}, err
		}
		item := shop.ProductOutPutItem{
			P:       p,
			Inven:   (p.InventoryID).Hex(),
			Cate:    cates,
			MediaID: avatar.ID.Hex(),
			URL:     avatar.URL,
		}
		list = append(list, item)
	}
	var shop1 models.Shop
	shop1, err1 := uc.repo.DetailShop(ctx, models.Scope{}, opt.ShopID)
	if err1 != nil {
		uc.l.Errorf(ctx, "shop.usecase.ListProduct.repo.DetailShop: %v", err1)

		return shop.GetProductOutput{}, err1
	}
	return shop.GetProductOutput{
		Products: list,
		Pag:      pag,
		Shop:     shop1,
	}, nil
}
func (uc implUsecase) IsValidProductID(ctx context.Context, productID primitive.ObjectID) bool {
	return uc.repo.IsValidProductID(ctx, productID)
}
