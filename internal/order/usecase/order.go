package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/internal/order/delivery/rabbitmq"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/internal/vouchers"
	"github.com/pt010104/api-golang/pkg/util"
)

func (uc implUseCase) CreateOrder(ctx context.Context, sc models.Scope, input order.CreateOrderInput) (models.Order, error) {
	checkoutModel, err := uc.repo.DetailCheckout(ctx, sc, input.CheckoutID)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.DetailCheckout", err)
		return models.Order{}, err
	}

	if checkoutModel.Status != models.CheckoutStatusPending {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.DetailCheckout", order.ErrCheckoutStatusInvalid)
		return models.Order{}, order.ErrCheckoutStatusInvalid
	}

	if checkoutModel.ExpiredAt.Before(util.Now()) {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.DetailCheckout", order.ErrCheckoutExpired)
		return models.Order{}, order.ErrCheckoutExpired
	}

	fmt.Println("checkoutModel.TotalPrice", checkoutModel.TotalPrice)

	var wg sync.WaitGroup
	var wgErr error
	totalPrice := checkoutModel.TotalPrice

	wg.Add(1)
	go func() {
		defer wg.Done()
		userModel, err := uc.userUC.GetModel(ctx, sc.UserID)
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateOrder.userUC.GetModel", err)
			wgErr = err
		}

		var existAddress bool
		for _, address := range userModel.Address {
			if address.ID.Hex() == input.AddressID {
				existAddress = true
				break
			}
		}
		if !existAddress {
			uc.l.Errorf(ctx, "order.usecase.CreateOrder.userUC.GetModel", order.ErrAddressNotFound)
			wgErr = order.ErrAddressNotFound
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if input.VoucherID != "" {
			var err error
			_, totalPrice, _, err = uc.voucherUC.ApplyVoucher(ctx, sc, vouchers.ApplyVoucherInput{
				ID:          input.VoucherID,
				OrderAmount: checkoutModel.TotalPrice,
			})
			if err != nil {
				uc.l.Errorf(ctx, "order.usecase.CreateOrder.voucherUC.ApplyVoucher", err)
				wgErr = err
			}
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return models.Order{}, wgErr
	}

	var orderModel models.Order

	wg.Add(1)
	go func() {
		defer wg.Done()
		checkoutModel, err = uc.repo.UpdateCheckout(ctx, sc, order.UpdateCheckoutOption{
			Model:  checkoutModel,
			Status: models.CheckoutStatusConfirmed,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.UpdateCheckout", err)
			wgErr = err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		orderModel, err = uc.repo.CreateOrder(ctx, sc, order.CreateOrderOption{
			CheckoutID:    checkoutModel.ID,
			Products:      checkoutModel.Products,
			PaymentMethod: input.PaymentMethod,
			AddressID:     input.AddressID,
			TotalPrice:    totalPrice,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateOrder.repo.CreateOrder", err)
			wgErr = err
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return models.Order{}, wgErr
	}

	err = uc.PublishOrder(ctx, sc, rabbitmq.OrderMessage{
		OrderID:    orderModel.ID.Hex(),
		CheckoutID: checkoutModel.ID.Hex(),
		UserID:     sc.UserID,
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.CreateOrder.PublishOrder", err)
		return models.Order{}, err
	}

	return orderModel, nil
}

func (uc implUseCase) DetailOrder(ctx context.Context, sc models.Scope, orderID string) (models.Order, error) {
	orderModel, err := uc.repo.DetailOrder(ctx, sc, orderID)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.DetailOrder.repo.DetailOrder", err)
		return models.Order{}, err
	}

	return orderModel, nil
}

func (uc implUseCase) ListOrder(ctx context.Context, sc models.Scope, input order.ListOrderInput) (order.ListOrderOutput, error) {
	orderModels, err := uc.repo.ListOrder(ctx, sc, order.ListOrderOption{
		Status: input.Status,
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.ListOrder.repo.ListOrder", err)
		return order.ListOrderOutput{}, err
	}

	if len(orderModels) == 0 {
		return order.ListOrderOutput{}, nil
	}

	productIDs := make([]string, 0)
	for _, order := range orderModels {
		for _, product := range order.Products {
			productIDs = append(productIDs, product.ID.Hex())
		}
	}

	productIDs = util.RemoveDuplicates(productIDs)

	var products []models.Product

	p, err := uc.shopUC.ListProduct(ctx, sc, shop.ListProductInput{
		GetProductFilter: shop.GetProductFilter{
			IDs: productIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.validateProducts.shopUC.ListProduct: %v", err)
		return order.ListOrderOutput{}, err
	}

	for _, product := range p.Products {
		products = append(products, product.P)
	}

	imageMap := make(map[string]string)
	for _, product := range p.Products {
		if len(product.Images) > 0 {
			imageMap[product.P.ID.Hex()] = product.Images[0].URL
		}
	}

	productMap := make(map[string]models.Product)
	for _, product := range products {
		productMap[product.ID.Hex()] = product
	}
	orderItems := make([]order.OrderItem, len(orderModels))
	for i, orderModel := range orderModels {
		fmt.Println("orderModel.ID: ", orderModel.ID.Hex())
		fmt.Println("orderModel.Products: ")
		util.PrintJson(orderModel.Products)
		productItems := make([]order.ProductItem, 0)
		for _, product := range orderModel.Products {
			fmt.Println("product.ID: ", product.ID.Hex())
			productItems = append(productItems, order.ProductItem{
				ProductID:   product.ID.Hex(),
				ProductName: productMap[product.ID.Hex()].Name,
				ImageURL:    imageMap[product.ID.Hex()],
				Price:       productMap[product.ID.Hex()].Price,
				Quantity:    product.Quantity,
			})
		}
		orderItems[i] = order.OrderItem{
			Order:      orderModel,
			Products:   productItems,
			TotalPrice: orderModel.TotalPrice,
		}

		fmt.Println("orderItems[i]: ")
		util.PrintJson(orderItems[i])
	}

	return order.ListOrderOutput{
		Orders: orderItems,
	}, nil
}

func (uc implUseCase) ListOrderShop(ctx context.Context, sc models.Scope, input order.ListOrderShopInput) (order.ListOrderShopOutput, error) {
	orderModels, err := uc.repo.ListOrderShop(ctx, sc, order.ListOrderOption{
		Status: input.Status,
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.ListOrderShop.repo.ListOrderShop", err)
		return order.ListOrderShopOutput{}, err
	}

	if len(orderModels) == 0 {
		return order.ListOrderShopOutput{}, nil
	}

	var (
		userModels    []models.User
		addressMerges []string
		wg            sync.WaitGroup
		mu            sync.Mutex
		wgErr         error
	)

	userDataChan := make(chan struct {
		user    models.User
		address string
	}, len(orderModels))

	for _, orderModel := range orderModels {
		wg.Add(1)
		go func(order models.Order) {
			defer wg.Done()

			userModel, err := uc.userUC.GetModel(ctx, order.UserID.Hex())
			if err != nil {
				uc.l.Errorf(ctx, "order.usecase.ListOrderShop.userUC.GetModel", err)
				mu.Lock()
				wgErr = err
				mu.Unlock()
				return
			}

			addressMerge := ""
			for _, address := range userModel.Address {
				if address.ID.Hex() == order.AddressID.Hex() {
					addressMerge = address.Street + ", " + address.District + ", " + address.City + ", " + address.Province
					break
				}
			}

			userDataChan <- struct {
				user    models.User
				address string
			}{
				user:    userModel,
				address: addressMerge,
			}
		}(orderModel)
	}

	go func() {
		wg.Wait()
		close(userDataChan)
	}()

	productIDsChan := make(chan string, 100)
	go func() {
		for _, order := range orderModels {
			for _, product := range order.Products {
				if product.ShopID.Hex() == sc.ShopID {
					productIDsChan <- product.ID.Hex()
				}
			}
		}
		close(productIDsChan)
	}()

	productIDsMap := make(map[string]struct{})
	for id := range productIDsChan {
		productIDsMap[id] = struct{}{}
	}

	productIDs := make([]string, 0, len(productIDsMap))
	for id := range productIDsMap {
		productIDs = append(productIDs, id)
	}

	for userData := range userDataChan {
		mu.Lock()
		userModels = append(userModels, userData.user)
		if userData.address != "" {
			addressMerges = append(addressMerges, userData.address)
		}
		mu.Unlock()
	}

	if wgErr != nil {
		return order.ListOrderShopOutput{}, wgErr
	}

	p, err := uc.shopUC.ListProduct(ctx, sc, shop.ListProductInput{
		GetProductFilter: shop.GetProductFilter{
			IDs: productIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.ListOrderShop.shopUC.ListProduct: %v", err)
		return order.ListOrderShopOutput{}, err
	}

	var imageMapMu, productMapMu sync.Mutex
	imageMap := make(map[string]string)
	productMap := make(map[string]models.Product)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for _, product := range p.Products {
			if len(product.Images) > 0 {
				imageMapMu.Lock()
				imageMap[product.P.ID.Hex()] = product.Images[0].URL
				imageMapMu.Unlock()
			}
		}
	}()

	go func() {
		defer wg.Done()
		for _, product := range p.Products {
			productMapMu.Lock()
			productMap[product.P.ID.Hex()] = product.P
			productMapMu.Unlock()
		}
	}()

	wg.Wait()

	orderItems := make([]order.OrderItem, 0, len(orderModels))
	for _, orderModel := range orderModels {
		productItems := make([]order.ProductItem, 0)
		for _, product := range orderModel.Products {
			if product.ShopID.Hex() == sc.ShopID {
				productItems = append(productItems, order.ProductItem{
					ProductID:   product.ID.Hex(),
					ProductName: productMap[product.ID.Hex()].Name,
					ImageURL:    imageMap[product.ID.Hex()],
					Price:       productMap[product.ID.Hex()].Price,
					Quantity:    product.Quantity,
				})
			}
		}
		orderItems = append(orderItems, order.OrderItem{
			Order:      orderModel,
			Products:   productItems,
			TotalPrice: orderModel.TotalPrice,
		})
	}

	return order.ListOrderShopOutput{
		Orders:        orderItems,
		AddressMerges: addressMerges,
		UserModels:    userModels,
	}, nil
}

func (uc implUseCase) UpdateOrder(ctx context.Context, sc models.Scope, input order.UpdateOrderInput) error {
	orderModel, err := uc.repo.DetailOrder(ctx, sc, input.OrderID)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.UpdateOrder.repo.DetailOrder", err)
		return err
	}

	err = uc.repo.UpdateOrder(ctx, sc, order.UpdateOrderOption{
		Model:  orderModel,
		Status: input.Status,
	})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.UpdateOrder.repo.UpdateOrder", err)
		return err
	}

	return nil
}
