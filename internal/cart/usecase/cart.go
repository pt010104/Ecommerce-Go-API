package usecase

import (
	"context"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (uc implUseCase) Create(sc models.Scope, ctx context.Context, input cart.CreateCartInput, inputItem cart.CreateCartItemInput) (models.Cart, error) {

	pidPrimitive, err := primitive.ObjectIDFromHex(inputItem.ProductID)
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Create.ObjecrIDFromHeX.inputItem.ProductID")
		return models.Cart{}, err
	}
	p, err1 := uc.shopUc.DetailProduct(ctx, sc, pidPrimitive)
	if err1 != nil {
		uc.l.Errorf(ctx, "cart.Usecase.DetailProduct")
		return models.Cart{}, err1
	}

	if (p.Inventory.StockLevel) < 1 {
		return models.Cart{}, cart.ErrNotEnoughQuantity
	}
	uidPrimitive, err := primitive.ObjectIDFromHex(sc.UserID)
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Create.ObjecrIDFromHeX.Sc.UserID")
		return models.Cart{}, err
	}
	shopidPrimitive, err := primitive.ObjectIDFromHex(input.ShopID)
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Create.ObjecrIDFromHeX.input.ShopID")
		return models.Cart{}, err
	}
	cart, err := uc.repo.Create(cart.CreateCartOption{
		UserID: uidPrimitive,
		ShopID: shopidPrimitive,
	},
		cart.CreateCartItemOption{
			ProductID: pidPrimitive,
			Quantity:  inputItem.Quantity,
		},
		ctx,
	)
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Repo.Create")
		return models.Cart{}, err
	}
	return cart, nil
}

func (uc implUseCase) Update(ctx context.Context, opt cart.UpdateCartOption) (models.Cart, error) {

	if len(opt.NewItemList) == 0 {
		return models.Cart{}, cart.ErrEmptyItemList
	}
	var ids []string
	for _, item := range opt.NewItemList {

		ids = append(ids, item.ProductID.Hex())
	}
	listProduct, err := uc.shopUc.ListProduct(ctx, models.Scope{}, shop.GetProductFilter{
		IDs: ids,
	})
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.Update.ListProduct", err)
		return models.Cart{}, err
	}
	if len(listProduct.Products) != len(opt.NewItemList) {
		return models.Cart{}, cart.ErrInvalidProductID
	}
	for _, item := range opt.NewItemList {

		if item.Quantity <= 0 {
			return models.Cart{}, cart.ErrInvalidQuantity
		}

		if !uc.shopUc.IsValidProductID(ctx, item.ProductID) {
			return models.Cart{}, cart.ErrInvalidProductID
		}
	}
	existingCart, err := uc.repo.Get(ctx, opt.ID)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return models.Cart{}, cart.ErrCartNotFound
		}
		return models.Cart{}, err
	}

	if existingCart.UserID != opt.UserID {
		return models.Cart{}, cart.ErrUserMismatch
	}

	updatedCart, err := uc.repo.Update(ctx, opt)
	if err != nil {
		return models.Cart{}, err
	}
	return updatedCart, nil
}
func (uc implUseCase) ListCart(sc models.Scope, ctx context.Context, opt cart.GetCartFilter) ([]models.Cart, error) {
	carts, err := uc.repo.ListCart(sc, ctx, opt)
	if err != nil {
		uc.l.Errorf(ctx, "cart.Usecase.ListCart", err)
		return nil, err
	}
	return carts, nil
}
