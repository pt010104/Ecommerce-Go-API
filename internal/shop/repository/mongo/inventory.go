package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (repo implRepo) DetailInventory(ctx context.Context, id primitive.ObjectID) (models.Inventory, error) {
	col := repo.getInventoryCollection()

	filter, err := repo.buildInventoryDetailQuery(ctx, id)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.DetailInventory.buildInventoryDetailQuery: %v", err)
		return models.Inventory{}, err
	}

	var u models.Inventory
	err = col.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.DetailInventory.FindOne: %v", err)
		return models.Inventory{}, err
	}
	return u, nil
}

func (repo implRepo) ListInventory(ctx context.Context, sc models.Scope, ids []primitive.ObjectID) ([]models.Inventory, error) {
	col := repo.getInventoryCollection()

	filter, err := repo.buildInventoryQuery(ctx, sc, ids)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListInventory.buildInventoryQuery: %v", err)
		return []models.Inventory{}, err
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListInventory.Find: %v", err)
		return []models.Inventory{}, err
	}
	defer cursor.Close(ctx)

	var inventories []models.Inventory
	err = cursor.All(ctx, &inventories)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListInventory.All: %v", err)
		return []models.Inventory{}, err
	}

	return inventories, nil
}

func (repo implRepo) UpdateInventory(ctx context.Context, sc models.Scope, opt shop.UpdateInventoryOption) (models.Inventory, error) {
	col := repo.getInventoryCollection()

	filter, err := repo.buildInventoryDetailQuery(ctx, opt.Model.ID)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.UpdateInventory.buildInventoryDetailQuery: %v", err)
		return models.Inventory{}, err
	}

	nm, update, err := repo.buildInventoryUpdateModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.UpdateInventory.buildInventoryUpdateModel: %v", err)
		return models.Inventory{}, err
	}

	var u models.Inventory
	err = col.FindOneAndUpdate(ctx, filter, update).Decode(&u)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.UpdateInventory.FindOneAndUpdate: %v", err)
		return models.Inventory{}, err
	}

	return nm, nil
}

func (repo implRepo) DeleteInventory(ctx context.Context, sc models.Scope, ids []primitive.ObjectID) error {
	col := repo.getInventoryCollection()

	filter, err := repo.buildInventoryQuery(ctx, sc, ids)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.DeleteInventory.buildInventoryQuery: %v", err)
		return err
	}

	_, err = col.DeleteMany(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.DeleteInventory.DeleteMany: %v", err)
		return err
	}

	return nil
}

func (repo implRepo) GetMostSoldInventory(ctx context.Context, sc models.Scope) ([]models.Inventory, error) {
	col := repo.getInventoryCollection()
	filter := bson.M{
		"sold_quantity": bson.M{"$gt": 0},
	}
	cursor, err := col.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "sold_quantity", Value: -1}, {Key: "_id", Value: -1}}))
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.GetMostSoldInventory.Find: %v", err)
		return []models.Inventory{}, err
	}
	defer cursor.Close(ctx)

	var inventories []models.Inventory
	err = cursor.All(ctx, &inventories)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.GetMostSoldInventory.All: %v", err)
		return []models.Inventory{}, err
	}

	return inventories, nil
}
