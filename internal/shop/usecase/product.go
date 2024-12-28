package usecase

import (
	"context"
	"sync"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/media"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/paginator"
	"github.com/pt010104/api-golang/pkg/util"
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

	_, err := uc.adminUC.ListCategories(ctx, sc, admins.GetCategoriesFilter{
		IDs: input.CategoryIDs,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.product.usecase.createproduct.validatecategoryids", err)
		return models.Product{}, models.Inventory{}, shop.ErrNonExistCategory
	}

	if len(input.MediaIDs) > 0 {

		_, err := uc.mediaUC.List(ctx, sc, media.ListInput{
			GetFilter: media.GetFilter{
				IDs: input.MediaIDs,
			},
		})
		if err != nil {
			uc.l.Errorf(ctx, "shop.product.usecase.createproduct.detailmedia", err)
			return models.Product{}, models.Inventory{}, err
		}
	}

	media_ids := mongo.ObjectIDsFromHexOrNil(input.MediaIDs)
	shopID, err := primitive.ObjectIDFromHex(sc.ShopID)
	if err != nil {
		uc.l.Errorf(ctx, "invalid ShopID format: %v", err)
		return models.Product{}, models.Inventory{}, err
	}

	p, err := uc.repo.CreateProduct(ctx, sc, shop.CreateProductOption{
		Name:        input.Name,
		Price:       input.Price,
		InventoryID: inven.ID,
		ShopID:      shopID,
		CategoryID:  mongo.ObjectIDsFromHexOrNil(input.CategoryIDs),
		MediaIDs:    media_ids,
		Alias:       util.BuildAlias(input.Name),
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
		wgErr         error

		medias []models.Media
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		u, err = uc.repo.Detailproduct(ctx, productID)
		if err != nil {
			uc.l.Errorf(ctx, "shop.product.usecase.detail.product", err)
			wgErr = err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = uc.repo.UpdateViewProduct(ctx, productID)
		if err != nil {
			uc.l.Errorf(ctx, "shop.product.usecase.detail.product", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()

	if wgErr != nil {
		return shop.DetailProductOutput{}, wgErr
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
		if len(u.MediaIDs) > 0 {
			medias, err = uc.mediaUC.List(ctx, sc, media.ListInput{
				GetFilter: media.GetFilter{
					IDs: mongo.HexFromObjectIDsOrNil(u.MediaIDs),
				},
			})
			if err != nil {
				errCh <- err
				return
			}
		}
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
		ID:           u.ID.Hex(),
		Name:         u.Name,
		CategoryName: categoryNames,
		Category:     category,
		Inventory:    inventory,
		Medias:       medias,
		View:         u.View,
		Shop:         shopDetail,
		Description:  u.Description,
		Price:        u.Price,
	}

	return output, nil
}

func (uc implUsecase) ListProduct(ctx context.Context, sc models.Scope, input shop.ListProductInput) (shop.ListProductOutput, error) {
	var (
		products    []models.Product
		s           models.Shop
		inventories []models.Inventory
		categories  []models.Category
	)

	err := validateProductInput(ctx, sc, input.GetProductFilter)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.ListProduct.validateProductInput: %v", err)
		return shop.ListProductOutput{}, err
	}

	var wg sync.WaitGroup
	var wgErr error
	wg.Add(2)

	input.IDs = util.RemoveDuplicates(input.IDs)

	go func() {
		defer wg.Done()
		var err error
		products, err = uc.repo.ListProduct(ctx, sc, shop.GetProductFilter{
			IDs:    input.IDs,
			ShopID: input.ShopID,
			Search: input.Search,
		})
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.ListProduct.repo.ListProduct: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		s, err = uc.repo.DetailShop(ctx, models.Scope{}, input.ShopID)
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

	var inventoryIDs []string
	categoryIDSet := make(map[string]struct{})
	for _, p := range products {
		inventoryIDs = append(inventoryIDs, p.InventoryID.Hex())
		for _, catID := range p.CategoryID {
			categoryIDSet[catID.Hex()] = struct{}{}
		}
	}

	var categoryIDs []string
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		inventoryIDs = util.RemoveDuplicates(inventoryIDs)
		inventories, err = uc.repo.ListInventory(ctx, models.Scope{}, mongo.ObjectIDsFromHexOrNil(inventoryIDs))
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.ListProduct.repo.ListInventory: %v", err)
			wgErr = err
			return
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		categories, err = uc.adminUC.ListCategories(ctx, models.Scope{}, admins.GetCategoriesFilter{
			IDs: categoryIDs,
		})
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.ListProduct.adminUC.ListCategories: %v", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()

	if wgErr != nil {
		return shop.ListProductOutput{}, wgErr
	}

	categoryMap := make(map[string]models.Category)
	for _, cate := range categories {
		categoryMap[cate.ID.Hex()] = cate
	}

	inventoryMap := make(map[primitive.ObjectID]models.Inventory)
	for _, inv := range inventories {
		inventoryMap[inv.ID] = inv
	}

	//map productID to its mediaIDs
	//map productID to its array of models.Media
	mediaMap := make(map[string][]models.Media)
	for _, p := range products {
		media, er := uc.mediaUC.List(ctx, models.Scope{}, media.ListInput{
			GetFilter: media.GetFilter{IDs: mongo.HexFromObjectIDsOrNil(p.MediaIDs)},
		})
		if er != nil {
			uc.l.Errorf(ctx, "shop.usecase.ListProduct: %v", er)
			return shop.ListProductOutput{}, er
		}
		mediaMap[p.ID.Hex()] = media

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
			P:         p,
			Inventory: inventoryMap[p.InventoryID],
			Cate:      cates,
			Images:    mediaMap[p.ID.Hex()],
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

func (uc implUsecase) GetProduct(ctx context.Context, sc models.Scope, input shop.GetProductInput) (shop.GetProductOutput, error) {
	var (
		s           []models.Product
		pag         paginator.Paginator
		categories  []models.Category
		inventories []models.Inventory
		shop1       models.Shop
		wg          sync.WaitGroup
		wgErr       error
		mediaMap    map[primitive.ObjectID][]models.Media
		allMediaIDs []string
	)

	opt := shop.GetProductOption{
		GetProductFilter: shop.GetProductFilter{
			ShopID:  input.ShopID,
			IDs:     input.IDs,
			CateIDs: input.CateIDs,
			Search:  input.Search,
		},
		PagQuery: input.PagQuery,
	}

	var err error
	s, pag, err = uc.repo.GetProduct(ctx, models.Scope{}, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.GetProduct.repo.GetProduct: %v", err)
		return shop.GetProductOutput{}, err
	}

	var inventoryIDs []string
	categoryIDSet := make(map[string]struct{})
	for _, p := range s {
		inventoryIDs = append(inventoryIDs, p.InventoryID.Hex())
		for _, catID := range p.CategoryID {
			categoryIDSet[catID.Hex()] = struct{}{}
		}
		allMediaIDs = append(allMediaIDs, mongo.HexFromObjectIDsOrNil(p.MediaIDs)...)
	}

	var categoryIDs []string
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	allMediaIDs = util.RemoveDuplicates(allMediaIDs)

	wg.Add(3)

	go func() {
		defer wg.Done()
		var err error
		categories, err = uc.adminUC.ListCategories(ctx, models.Scope{}, admins.GetCategoriesFilter{
			IDs: categoryIDs,
		})
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetProduct.adminUC.ListCategories: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		inventoryIDs = util.RemoveDuplicates(inventoryIDs)
		inventories, err = uc.repo.ListInventory(ctx, models.Scope{}, mongo.ObjectIDsFromHexOrNil(inventoryIDs))
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetProduct.repo.ListInventory: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		shop1, err = uc.repo.DetailShop(ctx, models.Scope{}, opt.ShopID)
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetProduct.repo.DetailShop: %v", err)
			wgErr = err
			return
		}
	}()

	wg.Wait()

	if wgErr != nil {
		return shop.GetProductOutput{}, wgErr
	}

	categoryMap := make(map[string]models.Category)
	for _, cate := range categories {
		categoryMap[cate.ID.Hex()] = cate
	}

	inventoryMap := make(map[primitive.ObjectID]models.Inventory)
	for _, inv := range inventories {
		inventoryMap[inv.ID] = inv
	}

	if len(allMediaIDs) > 0 {
		allMedia, err := uc.mediaUC.List(ctx, models.Scope{}, media.ListInput{
			GetFilter: media.GetFilter{IDs: allMediaIDs},
		})
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetProduct.mediaUC.List: %v", err)
			return shop.GetProductOutput{}, err
		}

		mediaMap = make(map[primitive.ObjectID][]models.Media)
		for _, p := range s {
			var productMedia []models.Media
			for _, m := range allMedia {
				for _, mediaID := range p.MediaIDs {
					if m.ID == mediaID {
						productMedia = append(productMedia, m)
					}
				}
			}
			mediaMap[p.ID] = productMedia
		}
	}

	var list []shop.ProductOutPutItem
	for _, p := range s {
		var cates []models.Category
		for _, catID := range p.CategoryID {
			if cate, ok := categoryMap[catID.Hex()]; ok {
				cates = append(cates, cate)
			}
		}

		item := shop.ProductOutPutItem{
			P:         p,
			Inventory: inventoryMap[p.InventoryID],
			Cate:      cates,
			Images:    mediaMap[p.ID],
		}
		list = append(list, item)
	}

	return shop.GetProductOutput{
		Products: list,
		Pag:      pag,
		Shop:     shop1,
	}, nil
}

func (uc implUsecase) GetAll(ctx context.Context, sc models.Scope, input shop.GetProductOption) (shop.GetAllProductOutput, error) {
	var (
		s           []models.Product
		pag         paginator.Paginator
		categories  []models.Category
		inventories []models.Inventory
		shops       []models.Shop
		wg          sync.WaitGroup
		wgErr       error
		mediaMap    map[primitive.ObjectID][]models.Media
		allMediaIDs []string
	)

	opt := shop.GetProductOption{
		GetProductFilter: input.GetProductFilter,
		PagQuery:         input.PagQuery,
	}

	s, pag, err := uc.repo.GetProduct(ctx, models.Scope{}, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.GetAll.repo.GetProduct: %v", err)
		return shop.GetAllProductOutput{}, err
	}

	var inventoryIDs []string
	categoryIDSet := make(map[string]struct{})
	shopIDSet := make(map[primitive.ObjectID]struct{})

	for _, p := range s {
		inventoryIDs = append(inventoryIDs, p.InventoryID.Hex())
		for _, catID := range p.CategoryID {
			categoryIDSet[catID.Hex()] = struct{}{}
		}
		shopIDSet[p.ShopID] = struct{}{}
		allMediaIDs = append(allMediaIDs, mongo.HexFromObjectIDsOrNil(p.MediaIDs)...)
	}

	var categoryIDs []string
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	var shopIDs []primitive.ObjectID
	for sid := range shopIDSet {
		shopIDs = append(shopIDs, sid)
	}

	allMediaIDs = util.RemoveDuplicates(allMediaIDs)

	wg.Add(4)
	go func() {
		defer wg.Done()
		var err error
		categories, err = uc.adminUC.ListCategories(ctx, models.Scope{}, admins.GetCategoriesFilter{
			IDs: categoryIDs,
		})
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetAll.adminUC.ListCategories: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		inventoryIDs = util.RemoveDuplicates(inventoryIDs)
		inventories, err = uc.repo.ListInventory(ctx, models.Scope{}, mongo.ObjectIDsFromHexOrNil(inventoryIDs))
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetAll.repo.ListInventory: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		shops, err = uc.repo.ListShop(ctx, models.Scope{}, shop.GetShopsFilter{
			IDs: mongo.HexFromObjectIDsOrNil(shopIDs),
		})
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.GetAll.repo.ListShops: %v", err)
			wgErr = err
			return
		}
	}()

	go func() {
		defer wg.Done()
		if len(allMediaIDs) > 0 {
			allMedia, err := uc.mediaUC.List(ctx, models.Scope{}, media.ListInput{
				GetFilter: media.GetFilter{IDs: allMediaIDs},
			})
			if err != nil {
				uc.l.Errorf(ctx, "shop.usecase.GetAll.mediaUC.List: %v", err)
				wgErr = err
				return
			}

			mediaMap = make(map[primitive.ObjectID][]models.Media)
			for _, p := range s {
				var productMedia []models.Media
				for _, m := range allMedia {
					for _, mediaID := range p.MediaIDs {
						if m.ID == mediaID {
							productMedia = append(productMedia, m)
						}
					}
				}
				mediaMap[p.ID] = productMedia
			}
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return shop.GetAllProductOutput{}, wgErr
	}

	categoryMap := make(map[string]models.Category)
	for _, cate := range categories {
		categoryMap[cate.ID.Hex()] = cate
	}

	inventoryMap := make(map[primitive.ObjectID]models.Inventory)
	for _, inv := range inventories {
		inventoryMap[inv.ID] = inv
	}

	shopMap := make(map[primitive.ObjectID]models.Shop)
	for _, sh := range shops {
		shopMap[sh.ID] = sh
	}

	var list []shop.GetAllProductItem
	for _, p := range s {
		var cates []models.Category
		for _, catID := range p.CategoryID {
			if cate, ok := categoryMap[catID.Hex()]; ok {
				cates = append(cates, cate)
			}
		}

		item := shop.GetAllProductItem{
			P:         p,
			Inventory: inventoryMap[p.InventoryID],
			Cate:      cates,
			Images:    mediaMap[p.ID],
			Shop:      shopMap[p.ShopID],
		}
		list = append(list, item)
	}

	return shop.GetAllProductOutput{
		Products: list,
		Pag:      pag,
	}, nil
}

func (uc implUsecase) UpdateProduct(ctx context.Context, sc models.Scope, input shop.UpdateProductOption) (models.Product, error) {
	shop1, err := uc.repo.DetailShop(ctx, models.Scope{}, sc.ShopID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.UpdateProduct.repoDetailShop: %v", err)
		return models.Product{}, err
	}
	user, err := uc.userUC.Detail(ctx, models.Scope{}, shop1.UserID.Hex())

	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.UpdateProduct.userUC.Detail: %v", err)
		return models.Product{}, err
	}
	if user.User.ID.Hex() != sc.UserID {

		return models.Product{}, shop.ErrNoPermissionToUpdate
	}
	p, err := uc.repo.UpdateProduct(ctx, sc, shop.UpdateProductOption{
		Name:       input.Name,
		ID:         input.ID,
		Price:      input.Price,
		CategoryID: input.CategoryID,
		MediaIDs:   input.MediaIDs,
		Alias:      util.BuildAlias(input.Name),
		Model:      input.Model,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.UpdateProduct.repoUpdateProduct: %v", err)
		return models.Product{}, err
	}

	p1, err1 := uc.repo.Detailproduct(ctx, p.ID)
	if err1 != nil {
		uc.l.Errorf(ctx, "shop.usecase.UpdateProduct.repoDetailProduct: %v", err1)
		return models.Product{}, err1
	}

	invenModel, err := uc.repo.DetailInventory(ctx, p1.InventoryID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.UpdateProduct.DetailInventory: %v", err1)
		return models.Product{}, err1
	}

	_, err = uc.repo.UpdateInventory(ctx, sc, shop.UpdateInventoryOption{
		Model:           invenModel,
		StockLevel:      &input.StockLevel,
		ReorderLevel:    &input.ReorderLevel,
		ReorderQuantity: &input.ReorderQuantity,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.UpdateProduct.UpdateInventory: %v", err1)
		return models.Product{}, err1
	}

	return p1, nil
}
