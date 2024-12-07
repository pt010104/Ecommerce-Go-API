package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	cartCollection = "carts"
)

func (repo implRepo) getCartCollection() mongo.Collection {
	return *repo.database.Collection(cartCollection)
}
func (repo implRepo) Create(opt cart.CreateCartOption, opt2 cart.CreateCartItemOption, ctx context.Context) (models.Cart, error) {
	col := repo.getCartCollection()
	newCart, err := repo.buildCartModel(opt, ctx)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Create.buildCartModel", err)
		return models.Cart{}, err
	}
	var items []models.CartItem
	item := models.CartItem{
		ProductID: opt2.ProductID,
		AddedAt:   time.Now(),
		Quantity:  opt2.Quantity,
	}
	items = append(items, item)
	newCart.Items = items
	_, err = col.InsertOne(ctx, newCart)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Create.InsertOne", err)
		return models.Cart{}, err
	}
	return newCart, nil
}
func (repo implRepo) Get(ctx context.Context, ID primitive.ObjectID) (models.Cart, error) {
	col := repo.getCartCollection()
	filter, err := repo.buildCartDetailQuery(ctx, ID)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Get.buildCartQuery", err)
		return models.Cart{}, err
	}
	var cart models.Cart
	err = col.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return models.Cart{}, err
		}

		repo.l.Errorf(ctx, "Cart.Repo.Get.FindOne", err)
		return models.Cart{}, err
	}
	return cart, nil
}
func (repo implRepo) Update(ctx context.Context, opt cart.UpdateCartOption) (models.Cart, error) {
	col := repo.getCartCollection()
	filter, err := repo.buildCartDetailQuery(ctx, opt.ID)
	//print filter
	repo.l.Debugf(ctx, "filter", filter)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Update.buildCartDetailQuery", err)
		return models.Cart{}, err
	}
	//print opt.NewItemList

	repo.l.Debugf(ctx, "opt.NewItemList", opt.NewItemList)
	update := bson.M{
		"$set": bson.M{
			"items": opt.NewItemList,
		},
	}
	opt.Model.ID = opt.ID
	opt.Model.UpdatedAt = time.Now()
	opt.Model.ShopID = opt.ShopID
	opt.Model.UserID = opt.UserID
	opt.Model.Items = opt.NewItemList

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Cart{}, err
		}
		repo.l.Errorf(ctx, "Cart.Repo.Update.FindOneAndUpdate", err)
		return models.Cart{}, err
	}

	return opt.Model, nil

}
func (repo implRepo) ListCart(sc models.Scope, ctx context.Context, opt cart.GetCartFilter) ([]models.Cart, error) {
	col := repo.getCartCollection()

	filter, err := repo.buildCartQuery(sc, opt, ctx)
	repo.l.Debugf(ctx, "filter", filter)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.buildCartQuery: %v", err)
		return []models.Cart{}, err
	}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.ListCart.Find: %v", err)
		return []models.Cart{}, err
	}
	defer cursor.Close(ctx)
	var carts []models.Cart
	for cursor.Next(ctx) {
		var cart models.Cart
		if err := cursor.Decode(&cart); err != nil {
			repo.l.Errorf(ctx, "Failed to decode cart: %v", err)
			continue
		}
		carts = append(carts, cart)
	}

	repo.l.Debugf(ctx, "carts", carts)
	return carts, nil
}
