package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/inventory"
	"github.com/pt010104/api-golang/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "inventories"
)

func (repo implRepo) getInventoryCollection() mongo.Collection {
	return *repo.database.Collection(userCollection)
}

func (repo implRepo) Create(ctx context.Context, sc models.Scope, opt inventory.CreateOption) (models.Inventory, error) {
	col := repo.getInventoryCollection()

	u, err := repo.buildInventoryModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.buldUserModel: %v", err)
		return models.Inventory{}, err

	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.InsertOne: %v", err)
		return models.Inventory{}, err
	}
	return u, nil
}
