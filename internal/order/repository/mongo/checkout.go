package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	checkoutCollection = "checkouts"
)

func (repo implRepo) getCheckoutCollection() mongo.Collection {
	return *repo.database.Collection(checkoutCollection)
}

func (repo implRepo) CreateCheckout(ctx context.Context, sc models.Scope, opt order.CreateCheckoutOption) (models.Checkout, error) {
	col := repo.getCheckoutCollection()
	newCheckout, err := repo.buildCheckoutModel(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.CreateCheckout.buildCheckoutModel", err)
		return models.Checkout{}, err
	}

	_, err = col.InsertOne(ctx, newCheckout)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.CreateCheckout.InsertOne", err)
		return models.Checkout{}, err
	}

	return newCheckout, nil
}

func (repo implRepo) DetailCheckout(ctx context.Context, sc models.Scope, checkoutID string) (models.Checkout, error) {
	col := repo.getCheckoutCollection()

	filter, err := repo.buildCheckoutDetailQuery(ctx, sc, checkoutID)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.DetailCheckout.buildCheckoutDetailQuery", err)
		return models.Checkout{}, err
	}

	var checkout models.Checkout
	err = col.FindOne(ctx, filter).Decode(&checkout)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.DetailCheckout.FindOne", err)
		return models.Checkout{}, err
	}

	return checkout, nil
}

func (repo implRepo) UpdateCheckout(ctx context.Context, sc models.Scope, opt order.UpdateCheckoutOption) (models.Checkout, error) {
	col := repo.getCheckoutCollection()

	filter, err := repo.buildCheckoutDetailQuery(ctx, sc, opt.Model.ID.Hex())
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.UpdateCheckout.buildCheckoutDetailQuery", err)
		return models.Checkout{}, err
	}

	nm, update := repo.buildCheckoutUpdate(ctx, opt)

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.UpdateCheckout.UpdateOne", err)
		return models.Checkout{}, err
	}

	return nm, nil
}
