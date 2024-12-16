package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	cartCollection = "carts"
)

func (repo implRepo) getCartCollection() mongo.Collection {
	return *repo.database.Collection(cartCollection)
}

func (repo implRepo) Create(ctx context.Context, sc models.Scope, opt cart.CreateCartOption) (models.Cart, error) {
	col := repo.getCartCollection()
	newCart, err := repo.buildCartModel(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Create.buildCartModel", err)
		return models.Cart{}, err
	}

	_, err = col.InsertOne(ctx, newCart)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Create.InsertOne", err)
		return models.Cart{}, err
	}

	return newCart, nil
}

func (repo implRepo) GetOne(ctx context.Context, sc models.Scope, opt cart.GetOneOption) (models.Cart, error) {
	col := repo.getCartCollection()
	filter, err := repo.buildCartQuery(ctx, sc, opt.CartFilter)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.GetOne.buildCartQuery", err)
		return models.Cart{}, err
	}
	var cart models.Cart
	err = col.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return models.Cart{}, err
		}

		repo.l.Errorf(ctx, "Cart.Repo.GetOne.FindOne", err)
		return models.Cart{}, err
	}
	return cart, nil
}

func (repo implRepo) Update(ctx context.Context, sc models.Scope, opt cart.UpdateCartOption) (models.Cart, error) {
	col := repo.getCartCollection()
	filter, err := repo.buildCartDetailQuery(ctx, sc, opt.Model.ID)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Update.buildCartDetailQuery", err)
		return models.Cart{}, err
	}

	nm, update, err := repo.buildCartUpdateModel(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Update.buildCartUpdateModel", err)
		return models.Cart{}, err
	}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Cart{}, err
		}
		repo.l.Errorf(ctx, "Cart.Repo.Update.FindOneAndUpdate", err)
		return models.Cart{}, err
	}

	return nm, nil

}
func (repo implRepo) ListCart(ctx context.Context, sc models.Scope, opt cart.ListOption) ([]models.Cart, error) {
	col := repo.getCartCollection()

	filter, err := repo.buildCartQuery(ctx, sc, opt.CartFilter)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.buildCartQuery: %v", err)
		return nil, err
	}

	var carts []models.Cart
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.ListCart.Find: %v", err)
		return nil, err
	}

	err = cursor.All(ctx, &carts)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.ListCart.All: %v", err)
		return nil, err
	}
	return carts, nil
}
