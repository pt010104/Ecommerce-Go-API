package usecase

import (
	"context"
	"math"
	"sync"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/internal/resources"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUseCase) ConsumeOrderMsg(ctx context.Context, sc models.Scope, input order.ConsumeOrderMsgInput) error {
	var wg sync.WaitGroup
	var wgErr error
	var us models.User
	var orderModel models.Order
	var productMap map[string]shop.ProductOutPutItem
	var shops []models.Shop
	var userShops []models.User
	shopIDs := make([]string, 0, len(orderModel.Products))

	orderModel, err := uc.repo.DetailOrder(ctx, sc, input.OrderID)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.repo.DetailOrder", err)
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		productIDs := make([]string, 0, len(orderModel.Products))
		for _, product := range orderModel.Products {
			productIDs = append(productIDs, product.ID.Hex())
		}

		listProductOutput, err := uc.shopUC.ListProduct(ctx, sc, shop.ListProductInput{
			GetProductFilter: shop.GetProductFilter{
				IDs: productIDs,
			},
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.shopUC.ListProduct", err)
			wgErr = err
			return
		}

		for _, product := range listProductOutput.Products {
			shopIDs = append(shopIDs, product.P.ShopID.Hex())
		}

		productMap = make(map[string]shop.ProductOutPutItem)
		for _, product := range listProductOutput.Products {
			productMap[product.P.ID.Hex()] = product
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		us, err = uc.userUC.GetModel(ctx, input.UserID)
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.userUC.GetModel", err)
			wgErr = err
			return
		}

	}()

	wg.Wait()
	if wgErr != nil {
		return wgErr
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		shops, err = uc.shopUC.ListShop(ctx, sc, shop.GetShopsFilter{
			IDs: shopIDs,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.shopUC.ListShop", err)
			wgErr = err
			return
		}

		userShopIDs := make([]string, 0, len(shops))
		for _, shop := range shops {
			userShopIDs = append(userShopIDs, shop.UserID.Hex())
		}

		userShops, err = uc.userUC.ListUsers(ctx, sc, user.ListUserInput{
			IDs: userShopIDs,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.userUC.List", err)
			wgErr = err
			return
		}
	}()

	for _, productOrder := range orderModel.Products {
		wg.Add(1)

		go func(productOrder models.OrderProduct) {
			defer wg.Done()
			product, ok := productMap[productOrder.ID.Hex()]
			if !ok {
				uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.shopUC.UpdateInventory", order.ErrProductNotFound)
				wgErr = order.ErrProductNotFound
				return
			}

			stockLevel := product.Inventory.StockLevel - uint(productOrder.Quantity)
			reservedLevel := product.Inventory.ReservedLevel - uint(productOrder.Quantity)

			_, err := uc.shopUC.UpdateInventory(ctx, sc, shop.UpdateInventoryInput{
				ID:            product.P.InventoryID,
				StockLevel:    util.ToPointer(stockLevel),
				ReservedLevel: reservedLevel,
				SoldQuantity:  product.Inventory.SoldQuantity + uint(productOrder.Quantity),
			})
			if err != nil {
				uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.shopUC.UpdateInventory", err)
				wgErr = err
				return
			}
		}(productOrder)
	}

	wg.Wait()
	if wgErr != nil {
		return wgErr
	}

	orderProducts := make([]resources.OrderProductEmail, 0)
	totalPrice := 0.0

	shopIDMap := make(map[string]string)
	for _, shop := range shops {
		shopIDMap[shop.UserID.Hex()] = shop.ID.Hex()
	}

	for _, productOrder := range orderModel.Products {
		product, ok := productMap[productOrder.ID.Hex()]
		if !ok {
			continue
		}

		shopInfo, err := uc.shopUC.Detail(ctx, models.Scope{}, product.P.ShopID.Hex())
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.shopUC.DetailShop", err)
			return err
		}
		subTotal := math.Round(float64(product.P.Price)*float64(productOrder.Quantity)*100) / 100
		totalPrice += subTotal

		orderProducts = append(orderProducts, resources.OrderProductEmail{
			ProductID:   product.P.ID.Hex(),
			ProductName: product.P.Name,
			Quantity:    productOrder.Quantity,
			Price:       product.P.Price,
			ShopName:    shopInfo.S.Name,
			SubTotal:    subTotal,
			ShopID:      shopInfo.S.ID.Hex(),
		})
	}

	baseEmailData := resources.OrderEmailData{
		OrderID:       orderModel.ID.Hex(),
		OrderDate:     orderModel.CreatedAt,
		PaymentMethod: orderModel.PaymentMethod,
		CustomerName:  us.Name,
		CustomerEmail: us.Email,
		Products:      orderProducts,
		TotalPrice:    totalPrice,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		content, err := resources.GenerateOrderEmail(false, baseEmailData)
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.emailUC.generateOrderEmail", err)
			wgErr = err
			return
		}

		err = uc.emailUC.SendEmail(ctx, resources.EmailData{
			To:      us.Email,
			Subject: "Order Confirmation",
			Content: content,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.emailUC.SendEmail", err)
			wgErr = err
			return
		}
	}()

	for _, userShop := range userShops {
		wg.Add(1)
		go func(userShop models.User) {
			defer wg.Done()
			shopProducts := uc.filterProductsByShop(orderProducts, shopIDMap[userShop.ID.Hex()])
			shopTotal := uc.calculateShopTotal(shopProducts)

			shopEmailData := baseEmailData
			shopEmailData.Products = shopProducts
			shopEmailData.ShopTotal = shopTotal

			content, err := resources.GenerateOrderEmail(true, shopEmailData)
			if err != nil {
				uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.generateOrderEmail", err)
				wgErr = err
				return
			}

			err = uc.emailUC.SendEmail(ctx, resources.EmailData{
				To:      userShop.Email,
				Subject: "Order Confirmation",
				Content: content,
			})
			if err != nil {
				uc.l.Errorf(ctx, "order.usecase.ConsumeOrderMsg.emailUC.SendEmail", err)
				wgErr = err
				return
			}
		}(userShop)
	}

	wg.Wait()
	if wgErr != nil {
		return wgErr
	}

	return nil
}
