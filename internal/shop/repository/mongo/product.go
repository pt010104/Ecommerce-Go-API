package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	productCollection = "products"
)

func (repo implRepo) getProductCollection() mongo.Collection {
	return *repo.database.Collection(productCollection)
}

func (repo implRepo) CreateProduct(ctx context.Context, sc models.Scope, opt shop.CreateProductOption) (models.Product, error) {
	colP := repo.getProductCollection()

	p, err := repo.buildProductModel(opt, ctx)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repo.product.build:", err)
		return models.Product{}, err
	}

	_, err = colP.InsertOne(ctx, p)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repo.product.insertonE:", err)
		return models.Product{}, err
	}

	return p, nil
}
