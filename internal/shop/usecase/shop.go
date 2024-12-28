package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUsecase) Create(ctx context.Context, sc models.Scope, input shop.CreateShop) (models.Shop, error) {
	if input.City == "" || input.Name == "" || input.Street == "" {
		return models.Shop{}, shop.ErrInvalidInput
	}

	_, err := uc.repo.DetailShop(ctx, sc, "")
	if err == nil {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", shop.ErrShopExist)
		return models.Shop{}, shop.ErrShopExist
	}

	if !util.IsValidPhone(input.Phone) {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", shop.ErrInvalidPhone)
		return models.Shop{}, shop.ErrInvalidPhone
	}

	opt := shop.CreateShopOption{
		Name:     input.Name,
		Alias:    util.BuildAlias(input.Name),
		City:     input.City,
		Street:   input.Street,
		District: input.District,
		Phone:    input.Phone,
	}

	sh, err := uc.repo.CreateShop(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Create: %v", err)
		return models.Shop{}, err
	}

	return sh, nil
}

func (uc implUsecase) Get(ctx context.Context, sc models.Scope, input shop.GetShopInput) (shop.GetShopOutput, error) {
	opt := shop.GetOption{
		GetShopsFilter: input.GetShopsFilter,
		PagQuery:       input.PagQuery,
	}
	s, pag, err := uc.repo.GetShop(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Get: %v", err)
		return shop.GetShopOutput{}, err
	}

	var shop1 []shop.Shop_obj

	for _, v := range s {

		avatar, err := uc.userUC.Detail(ctx, models.Scope{}, v.UserID.Hex())
		if err != nil {
			uc.l.Errorf(ctx, "shop.usecase.Get: %v", err)
			return shop.GetShopOutput{}, err
		}

		shop1 = append(shop1, shop.Shop_obj{
			Shop: v,
			Avatar: shop.Avatar_obj{
				MediaID: avatar.User.MediaID.Hex(),
				URL:     avatar.Avatar_URL,
			},
		})

	}

	return shop.GetShopOutput{
		Shops: shop1,
		Pag:   pag,
	}, nil
}

func (uc implUsecase) Detail(ctx context.Context, sc models.Scope, id string) (shop.DetailShopOutput, error) {
	s, err := uc.repo.DetailShop(ctx, sc, id)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Detail: %v", err)
		return shop.DetailShopOutput{}, err
	}

	avatar, err := uc.userUC.Detail(ctx, sc, s.UserID.Hex())
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Detail: %v", err)
		return shop.DetailShopOutput{}, err
	}
	return shop.DetailShopOutput{
		S:       s,
		MediaID: avatar.User.MediaID.Hex(),
		URL:     avatar.Avatar_URL,
	}, nil
}

func (uc implUsecase) Delete(ctx context.Context, sc models.Scope) error {
	err := uc.repo.DeleteShop(ctx, sc)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Delete.Repodele", err)
		return shop.ErrShopDoesNotExist
	}

	return nil
}

func (uc implUsecase) Update(ctx context.Context, sc models.Scope, input shop.UpdateInput) ([]models.Shop, error) {
	var ids []string
	if input.ShopID != "" {
		ids = append(ids, input.ShopID)
	} else if len(input.ShopIDs) > 0 {
		ids = input.ShopIDs
	}

	ids = util.RemoveDuplicates(ids)

	ss, err := uc.repo.ListShop(ctx, sc, shop.GetShopsFilter{
		IDs: ids,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.update.repo.detail:", err)
		return []models.Shop{}, err
	}

	var shops []models.Shop

	if len(ss) > 0 {
		for _, s := range ss {
			var wgUpdate sync.WaitGroup
			var wgErrUpdate error
			var muUpdate sync.Mutex
			var ns models.Shop

			wgUpdate.Add(1)
			go func(s models.Shop) {
				defer wgUpdate.Done()
				ns, err = uc.repo.UpdateShop(ctx, sc, shop.UpdateOption{
					Model:      s,
					Name:       input.Name,
					Alias:      util.BuildAlias(input.Name),
					Phone:      input.Phone,
					City:       input.City,
					District:   input.District,
					Street:     input.Street,
					IsVerified: input.IsVerified,
				})
				if err != nil {
					uc.l.Errorf(ctx, "shop.usecase.update.repo.update:", err)
					wgErrUpdate = err
					return
				}

			}(s)

			if wgErrUpdate != nil {
				uc.l.Errorf(ctx, "shop.usecase.update.repo.update:", wgErrUpdate)
				return []models.Shop{}, wgErrUpdate
			}

			wgUpdate.Wait()
			shops = append(shops, ns)
			muUpdate.Lock()

		}

	}

	return shops, nil
}

func (uc implUsecase) ListShop(ctx context.Context, sc models.Scope, opt shop.GetShopsFilter) ([]models.Shop, error) {
	s, err := uc.repo.ListShop(ctx, sc, opt)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.ListShop: %v", err)
		return []models.Shop{}, err
	}

	return s, nil
}

func (uc implUsecase) GetIDByUserID(ctx context.Context, sc models.Scope, userID string) (string, error) {
	id, err := uc.repo.GetShopIDByUserID(ctx, sc, userID)
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.GetIDByUserID: %v", err)
		return "", err
	}

	return id, nil
}

func (uc implUsecase) Report(ctx context.Context, sc models.Scope) (shop.ReportOutput, error) {
	var (
		mostViewedProducts                           []models.Product
		mostSoldInventories                          []models.Inventory
		mostViewTrend                                []models.Product
		wg                                           sync.WaitGroup
		mu                                           sync.Mutex
		errMostViewed, errMostSold, errMostViewTrend error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		var err error
		mostViewedProducts, err = uc.repo.GetMostViewedProducts(ctx, sc)
		if err != nil {
			mu.Lock()
			errMostViewed = err
			mu.Unlock()
			uc.l.Errorf(ctx, "shop.usecase.Report.GetMostViewedProducts: %v", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		mostSoldInventories, err = uc.repo.GetMostSoldInventory(ctx, sc)
		if err != nil {
			mu.Lock()
			errMostSold = err
			mu.Unlock()
			uc.l.Errorf(ctx, "shop.usecase.Report.GetMostSoldInventory: %v", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		mostViewTrend, err = uc.repo.GetMostViewTrend(ctx, sc)
		if err != nil {
			mu.Lock()
			errMostViewTrend = err
			mu.Unlock()
			uc.l.Errorf(ctx, "shop.usecase.Report.GetMostViewTrend: %v", err)
			return
		}
	}()

	wg.Wait()

	if errMostViewed != nil {
		return shop.ReportOutput{}, errMostViewed
	}

	if errMostSold != nil {
		return shop.ReportOutput{}, errMostSold
	}

	if errMostViewTrend != nil {
		return shop.ReportOutput{}, errMostViewTrend
	}

	invenMap := make(map[string]int)
	var inventoryIDs []string
	for _, v := range mostSoldInventories {
		inventoryIDs = append(inventoryIDs, v.ID.Hex())
		invenMap[v.ID.Hex()] = int(v.SoldQuantity)
	}

	fmt.Println(inventoryIDs)

	if len(inventoryIDs) == 0 {
		return shop.ReportOutput{
			MostViewedProducts: mostViewedProducts,
			MostSoldProducts:   []shop.MostSolProductItem{},
		}, nil
	}

	products, err := uc.repo.ListProduct(ctx, sc, shop.GetProductFilter{
		InventoryIDs: inventoryIDs,
		ShopID:       sc.ShopID,
	})
	if err != nil {
		uc.l.Errorf(ctx, "shop.usecase.Report.ListProduct: %v", err)
		return shop.ReportOutput{}, err
	}

	var mostSoldProducts []shop.MostSolProductItem

	for _, v := range products {
		mostSoldProducts = append(mostSoldProducts, shop.MostSolProductItem{
			ProductID:   v.ID.Hex(),
			ProductName: v.Name,
			Sold:        invenMap[v.InventoryID.Hex()],
		})
	}

	return shop.ReportOutput{
		MostViewedProducts: mostViewedProducts,
		MostSoldProducts:   mostSoldProducts,
		MostViewTrend:      mostViewTrend,
	}, nil
}
