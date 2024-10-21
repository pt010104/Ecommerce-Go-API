package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "shop_inventories"
)

func (repo implRepo) getInventoryCollection() mongo.Collection {
	return *repo.database.Collection(userCollection)
}

func (repo implRepo) CreateInventory(ctx context.Context, sc models.Scope, opt shop.CreateInventoryOption) (models.Inventory, error) {
	col := repo.getInventoryCollection()

	u, err := repo.buildInventoryModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.CreateInventory.buildInventoryModel: %v", err)
		return models.Inventory{}, err

	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.CreateInventory.InsertOne: %v", err)
		return models.Inventory{}, err
	}
	return u, nil
}
