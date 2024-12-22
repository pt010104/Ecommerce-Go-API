package usecase

import (
	"context"
	"slices"
	"sync"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/checkout"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUseCase) validateProducts(ctx context.Context, sc models.Scope, productIDs []string) ([]models.Product, []string, error) {
	p, err := uc.shopUC.ListProduct(ctx, sc, shop.ListProductInput{
		GetProductFilter: shop.GetProductFilter{
			IDs: productIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "checkout.usecase.validateProducts.shopUC.ListProduct: %v", err)
		return nil, nil, err
	}

	if len(p.Products) != len(productIDs) {
		uc.l.Errorf(ctx, "checkout.usecase.validateProducts: %v", checkout.ErrProductNotFound)
		return nil, nil, checkout.ErrProductNotFound
	}

	var products []models.Product
	for _, product := range p.Products {
		products = append(products, product.P)
	}

	var image_urls []string
	for _, product := range p.Products {
		image_urls = append(image_urls, product.Images[0].URL)
	}

	return products, image_urls, nil
}

func (uc implUseCase) validateCart(ctx context.Context, sc models.Scope, productIDs []string) ([]primitive.ObjectID, map[string]int, error) {
	carts, err := uc.cartUC.GetCart(ctx, sc, cart.GetOption{})
	if err != nil {
		uc.l.Errorf(ctx, "checkout.usecase.validateCart.cartUC.GetCart: %v", err)
		return nil, nil, err
	}

	var cartIDs []primitive.ObjectID
	var cartProductIDs []string
	productQuantityMap := make(map[string]int)

	for _, CartOutPuts := range carts.CartOutPut {
		for _, p := range CartOutPuts.Products {
			cartProductIDs = append(cartProductIDs, p.ProductID)
			productQuantityMap[p.ProductID] = p.Quantity
		}
		cartIDs = append(cartIDs, CartOutPuts.Cart.ID)
	}

	for _, pid := range productIDs {
		if !slices.Contains(cartProductIDs, pid) {
			uc.l.Errorf(ctx, "checkout.usecase.validateCart: %v", checkout.ErrProductNotFoundInCart)
			return nil, nil, checkout.ErrProductNotFoundInCart
		}
	}

	return cartIDs, productQuantityMap, nil
}

func (uc implUseCase) validateStock(ctx context.Context, sc models.Scope, products []models.Product, productQuantityMap map[string]int) ([]models.Inventory, error) {
	inventoryIDs := make([]primitive.ObjectID, len(products))
	for _, product := range products {
		inventoryIDs = append(inventoryIDs, product.InventoryID)
	}

	invens, err := uc.shopUC.ListInventory(ctx, sc, inventoryIDs)
	if err != nil {
		uc.l.Errorf(ctx, "checkout.usecase.validateStock.shopUC.ListInventory: %v", err)
		return nil, err
	}

	for _, inven := range invens {
		if inven.StockLevel < uint(productQuantityMap[inven.ID.Hex()]) {
			uc.l.Errorf(ctx, "checkout.usecase.validateStock: %v", checkout.ErrProductNotEnoughStock)
			return nil, checkout.ErrProductNotEnoughStock
		}
	}

	return invens, nil
}

func (uc implUseCase) updateReservedLevel(ctx context.Context, sc models.Scope, invens []models.Inventory, inventoryIDs []primitive.ObjectID, productQuantityMap map[string]int, products []models.Product) error {
	var reservedLevelMap = make(map[string]uint)
	for _, inventory := range invens {
		reservedLevelMap[inventory.ID.Hex()] = inventory.ReservedLevel
	}

	for _, product := range products {
		reservedLevelMap[product.InventoryID.Hex()] += uint(productQuantityMap[product.ID.Hex()])
	}

	var wg sync.WaitGroup
	var wgErr error
	wg.Add(len(inventoryIDs))

	for _, inventoryID := range inventoryIDs {
		go func(inventoryID primitive.ObjectID) {
			defer wg.Done()
			_, err := uc.shopUC.UpdateInventory(ctx, sc, shop.UpdateInventoryInput{
				ID:            inventoryID,
				ReservedLevel: reservedLevelMap[inventoryID.Hex()],
			})
			if err != nil {
				wgErr = err
			}
		}(inventoryID)
	}

	wg.Wait()
	if wgErr != nil {
		uc.l.Errorf(ctx, "checkout.usecase.updateReservedLevel: %v", wgErr)
		return wgErr
	}

	return nil
}

func (uc implUseCase) calculateTotalPrices(ctx context.Context, sc models.Scope, products []models.Product, productQuantityMap map[string]int) ([]models.Shop, map[string]float64, float64, error) {
	shopIDs := make([]primitive.ObjectID, len(products))
	for _, product := range products {
		shopIDs = append(shopIDs, product.ShopID)
	}

	shops, err := uc.shopUC.ListShop(ctx, sc, shop.GetShopsFilter{
		IDs: mongo.HexFromObjectIDsOrNil(shopIDs),
	})
	if err != nil {
		uc.l.Errorf(ctx, "checkout.usecase.calculateTotalPrices.shopUC.ListShop: %v", err)
		return nil, nil, 0, err
	}

	totalPricesByShop, err := uc.calculateTotalPricesByShop(ctx, sc, products, productQuantityMap)
	if err != nil {
		uc.l.Errorf(ctx, "checkout.usecase.calculateTotalPrices: %v", err)
		return nil, nil, 0, err
	}

	var totalPrice float64
	for _, price := range totalPricesByShop {
		totalPrice += price
	}

	return shops, totalPricesByShop, totalPrice, nil
}

func (uc implUseCase) calculateTotalPricesByShop(ctx context.Context, sc models.Scope, products []models.Product, productQuantityMap map[string]int) (map[string]float64, error) {
	totalPricesByShop := make(map[string]float64)

	for _, product := range products {
		totalPricesByShop[product.ShopID.Hex()] += float64(product.Price) * float64(productQuantityMap[product.ID.Hex()])
	}

	return totalPricesByShop, nil
}
