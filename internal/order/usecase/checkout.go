package usecase

import (
	"context"
	"sync"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUseCase) CreateCheckout(ctx context.Context, sc models.Scope, productIDs []string) (order.CreateCheckoutOutput, error) {
	var wg sync.WaitGroup
	var wgErr error
	wg.Add(2)

	var products []models.Product
	var image_urls []string
	var productQuantityMap map[string]int

	go func() {
		defer wg.Done()
		var err error
		products, image_urls, err = uc.validateProducts(ctx, sc, productIDs)
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateCheckout.validateProducts: %v", err)
			wgErr = err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		_, productQuantityMap, err = uc.validateCart(ctx, sc, productIDs)
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateCheckout.validateCart: %v", err)
			wgErr = err
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return order.CreateCheckoutOutput{}, wgErr
	}

	invens, err := uc.validateStock(ctx, sc, products, productQuantityMap)
	if err != nil {
		return order.CreateCheckoutOutput{}, err
	}

	inventoryIDs := make([]primitive.ObjectID, 0, len(invens))
	for _, inven := range invens {
		inventoryIDs = append(inventoryIDs, inven.ID)
	}

	err = uc.updateReservedLevel(ctx, sc, invens, inventoryIDs, productQuantityMap, products)
	if err != nil {
		return order.CreateCheckoutOutput{}, err
	}

	var checkoutModel models.Checkout
	wg.Add(1)
	go func() {
		defer wg.Done()
		var productCheckouts []models.OrderProduct
		for _, product := range products {
			productCheckouts = append(productCheckouts, models.OrderProduct{
				ID:       product.ID,
				Quantity: productQuantityMap[product.ID.Hex()],
			})
		}
		checkoutModel, err = uc.repo.CreateCheckout(ctx, sc, order.CreateCheckoutOption{
			Products: productCheckouts,
		})
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateCheckout.repo.CreateCheckout: %v", err)
			wgErr = err
		}
	}()

	var shops []models.Shop
	var totalPricesByShop map[string]float64
	var totalPrice float64
	wg.Add(1)
	go func() {
		defer wg.Done()
		shops, totalPricesByShop, totalPrice, err = uc.calculateTotalPrices(ctx, sc, products, productQuantityMap)
		if err != nil {
			uc.l.Errorf(ctx, "order.usecase.CreateCheckout.calculateTotalPrices: %v", err)
			wgErr = err
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return order.CreateCheckoutOutput{}, wgErr
	}

	return order.CreateCheckoutOutput{
		CheckoutID:       checkoutModel.ID.Hex(),
		ExpiredAt:        checkoutModel.ExpiredAt,
		TotalPriceByShop: totalPricesByShop,
		TotalPrice:       totalPrice,
		Products:         products,
		QuantityMap:      productQuantityMap,
		Shops:            shops,
		ImageURLs:        image_urls,
	}, nil
}
