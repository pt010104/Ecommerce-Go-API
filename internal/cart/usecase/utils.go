package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
)

func (uc implUseCase) validateCartItem(ctx context.Context, sc models.Scope, input []cart.CartItemInput) error {
	var productMap = make(map[string]struct{})
	for _, item := range input {
		if _, ok := productMap[item.ProductID]; ok {
			uc.l.Errorf(ctx, "cart.Usecase.validateCartItem: %v", cart.ErrInvalidCartItem)
			return cart.ErrInvalidCartItem
		}
		if item.ProductID == "" || item.Quantity < 0 || !mongo.IsObjectID(item.ProductID) {
			uc.l.Errorf(ctx, "cart.Usecase.validateCartItem: %v", cart.ErrInvalidCartItem)
			return cart.ErrInvalidCartItem
		}
		productMap[item.ProductID] = struct{}{}
	}

	return nil
}

func (uc implUseCase) checkStock(ctx context.Context, sc models.Scope, inventory models.Inventory, quantity int) error {
	if inventory.StockLevel < uint(quantity) {
		uc.l.Errorf(ctx, "cart.Usecase.checkStock: %v", cart.ErrNotEnoughStock)
		return cart.ErrNotEnoughStock
	}
	return nil
}

func (uc implUseCase) getDataCartItems(ctx context.Context, sc models.Scope, input []cart.CartItemInput) (getDataOutput, error) {
	var productIDs []string
	var quantity []int
	for _, item := range input {
		productIDs = append(productIDs, item.ProductID)
		quantity = append(quantity, item.Quantity)
	}

	productOutput, err := uc.shopUc.ListProduct(ctx, sc, shop.ListProductInput{
		GetProductFilter: shop.GetProductFilter{
			IDs: productIDs,
		},
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.getDataCartItems.shopUc.ListProduct: %v", err)
		return getDataOutput{}, err
	}

	if len(productOutput.Products) != len(productIDs) {
		uc.l.Errorf(ctx, "cart.Usecase.getDataCartItems: %v", cart.ErrInvalidCartItem)
		return getDataOutput{}, cart.ErrInvalidCartItem
	}

	productMap := make(map[string]shop.ProductOutPutItem)
	for _, product := range productOutput.Products {
		productMap[product.P.ID.Hex()] = product
	}

	for _, item := range input {
		if err := uc.checkStock(ctx, sc, productMap[item.ProductID].Inventory, item.Quantity); err != nil {
			uc.l.Errorf(ctx, "cart.Usecase.getDataCartItems.checkStock: %v", err)
			return getDataOutput{}, err
		}
	}

	var shopIDs []string
	for _, product := range productOutput.Products {
		shopIDs = append(shopIDs, product.P.ShopID.Hex())
	}

	shops, err := uc.shopUc.ListShop(ctx, sc, shop.GetShopsFilter{
		IDs: shopIDs,
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.getDataCartItems.shopUc.ListShop: %v", err)
		return getDataOutput{}, err
	}

	if len(shops) != len(shopIDs) && len(shops) != len(productOutput.Products) {
		uc.l.Errorf(ctx, "cart.Usecase.getDataCartItems: %v", cart.ErrInvalidCartItem)
		return getDataOutput{}, cart.ErrInvalidCartItem
	}

	var cartItems []models.CartItem
	for i, product := range productOutput.Products {
		cartItems = append(cartItems, models.CartItem{
			ProductID: product.P.ID,
			Quantity:  quantity[i],
		})
	}

	return getDataOutput{
		Shops:      shops,
		CartItems:  cartItems,
		ShopIDs:    shopIDs,
		ProductMap: productMap,
	}, nil
}
