package usecase

import (
	"context"
	"slices"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/internal/resources"
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
		uc.l.Errorf(ctx, "order.usecase.validateProducts.shopUC.ListProduct: %v", err)
		return nil, nil, err
	}

	if len(p.Products) != len(productIDs) {
		uc.l.Errorf(ctx, "order.usecase.validateProducts: %v", order.ErrProductNotFound)
		return nil, nil, order.ErrProductNotFound
	}

	var products []models.Product
	for _, product := range p.Products {
		products = append(products, product.P)
	}

	image_urls := make([]string, len(p.Products))
	for i, product := range p.Products {
		if len(product.Images) > 0 {
			image_urls[i] = product.Images[0].URL
		}
	}

	return products, image_urls, nil
}

func (uc implUseCase) validateCart(ctx context.Context, sc models.Scope, productIDs []string) ([]primitive.ObjectID, map[string]int, error) {
	carts, err := uc.cartUC.GetCart(ctx, sc, cart.GetOption{})
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.validateCart.cartUC.GetCart: %v", err)
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
			uc.l.Errorf(ctx, "order.usecase.validateCart: %v", order.ErrProductNotFoundInCart)
			return nil, nil, order.ErrProductNotFoundInCart
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
		uc.l.Errorf(ctx, "order.usecase.validateStock.shopUC.ListInventory: %v", err)
		return nil, err
	}

	for _, inven := range invens {
		if inven.StockLevel-inven.ReservedLevel < uint(productQuantityMap[inven.ID.Hex()]) {
			uc.l.Errorf(ctx, "order.usecase.validateStock: %v", order.ErrProductNotEnoughStock)
			return nil, order.ErrProductNotEnoughStock
		}
	}

	return invens, nil
}

func (uc implUseCase) updateReservedLevel(ctx context.Context, sc models.Scope, invens []models.Inventory, inventoryIDs []primitive.ObjectID, productQuantityMap map[string]int, products []models.Product) error {
	const maxRetries = 5
	for _, inventory := range invens {
		retryCount := 0
		for retryCount < maxRetries {
			var currentVersion int
			currentVersion, err := uc.redisRepo.GetVersionInventory(ctx, inventory.ID)
			if err != nil {
				if err.Error() == order.ErrRedisNotFound.Error() {
					currentVersion = 0
				} else {
					uc.l.Errorf(ctx, "order.usecase.updateReservedLevel.redisRepo.GetVersionInventory: %v", err)
					return err
				}
			}

			newReservedLevel := inventory.ReservedLevel
			for _, product := range products {
				if product.InventoryID == inventory.ID {
					newReservedLevel += uint(productQuantityMap[product.ID.Hex()])
				}
			}

			_, err = uc.shopUC.UpdateInventory(ctx, sc, shop.UpdateInventoryInput{
				ID:            inventory.ID,
				ReservedLevel: newReservedLevel,
			})
			if err != nil {
				uc.l.Errorf(ctx, "order.usecase.updateReservedLevel.shopUC.UpdateInventory: %v", err)
				return err
			}

			newVersion, err := uc.redisRepo.IncrementVersionInventory(ctx, inventory.ID)
			if err != nil {
				uc.l.Errorf(ctx, "order.usecase.updateReservedLevel.redisRepo.IncrementVersionInventory: %v", err)
				return err
			}

			if newVersion != currentVersion+1 {
				retryCount++
				freshInventory, err := uc.shopUC.DetailInventory(ctx, inventory.ID)
				if err != nil {
					uc.l.Errorf(ctx, "order.usecase.updateReservedLevel.shopUC.DetailInventory: %v", err)
					return err
				}

				if freshInventory.StockLevel-freshInventory.ReservedLevel < uint(productQuantityMap[inventory.ID.Hex()]) {
					uc.l.Errorf(ctx, "order.usecase.updateReservedLevel: %v", order.ErrProductNotEnoughStock)
					return order.ErrProductNotEnoughStock
				}

				inventory = freshInventory
				continue
			}

			break
		}

		if retryCount == maxRetries {
			uc.l.Errorf(ctx, "order.usecase.updateReservedLevel: too many retries due to concurrent modifications")
			return order.ErrTooManyRetries
		}
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
		uc.l.Errorf(ctx, "order.usecase.calculateTotalPrices.shopUC.ListShop: %v", err)
		return nil, nil, 0, err
	}

	totalPricesByShop, err := uc.calculateTotalPricesByShop(ctx, sc, products, productQuantityMap)
	if err != nil {
		uc.l.Errorf(ctx, "order.usecase.calculateTotalPrices: %v", err)
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

func (uc implUseCase) filterProductsByShop(products []resources.OrderProductEmail, shopID string) []resources.OrderProductEmail {
	var filtered []resources.OrderProductEmail
	for _, p := range products {
		if p.ShopID == shopID {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func (uc implUseCase) calculateShopTotal(products []resources.OrderProductEmail) float64 {
	total := 0.0
	for _, p := range products {
		total += p.SubTotal
	}
	return total
}
