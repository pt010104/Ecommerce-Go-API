package usecase

import (
	"context"
	"sync"

	"github.com/pt010104/api-golang/internal/checkout"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUseCase) Create(ctx context.Context, sc models.Scope, productIDs []string) (checkout.CreateOutput, error) {
	var wg sync.WaitGroup
	var wgErr error
	wg.Add(2)

	var products []models.Product
	var image_urls []string
	var cartIDs []primitive.ObjectID
	var productQuantityMap map[string]int

	go func() {
		defer wg.Done()
		var err error
		products, image_urls, err = uc.validateProducts(ctx, sc, productIDs)
		if err != nil {
			uc.l.Errorf(ctx, "checkout.usecase.Create.validateProducts: %v", err)
			wgErr = err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		cartIDs, productQuantityMap, err = uc.validateCart(ctx, sc, productIDs)
		if err != nil {
			uc.l.Errorf(ctx, "checkout.usecase.Create.validateCart: %v", err)
			wgErr = err
		}
	}()

	wg.Wait()
	if wgErr != nil {
		return checkout.CreateOutput{}, wgErr
	}

	invens, err := uc.validateStock(ctx, sc, products, productQuantityMap)
	if err != nil {
		return checkout.CreateOutput{}, err
	}

	checkoutModel, err := uc.repo.Create(ctx, sc, checkout.CreateOption{
		CartIDs: cartIDs,
	})
	if err != nil {
		uc.l.Errorf(ctx, "checkout.usecase.Create.repo.Create: %v", err)
		return checkout.CreateOutput{}, err
	}

	inventoryIDs := make([]primitive.ObjectID, 0, len(invens))
	for _, inven := range invens {
		inventoryIDs = append(inventoryIDs, inven.ID)
	}

	err = uc.updateReservedLevel(ctx, sc, invens, inventoryIDs, productQuantityMap, products)
	if err != nil {
		return checkout.CreateOutput{}, err
	}

	shops, totalPricesByShop, totalPrice, err := uc.calculateTotalPrices(ctx, sc, products, productQuantityMap)
	if err != nil {
		return checkout.CreateOutput{}, err
	}

	return checkout.CreateOutput{
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
