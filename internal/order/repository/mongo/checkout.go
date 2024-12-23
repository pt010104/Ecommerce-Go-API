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
		repo.l.Errorf(ctx, "Checkout.Repo.CreateCheckout.buildCheckoutModel", err)
		return models.Checkout{}, err
	}

	_, err = col.InsertOne(ctx, newCheckout)
	if err != nil {
		repo.l.Errorf(ctx, "Checkout.Repo.CreateCheckout.InsertOne", err)
		return models.Checkout{}, err
	}

	return newCheckout, nil
}
