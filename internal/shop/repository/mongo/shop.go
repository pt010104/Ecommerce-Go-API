package mongo

import (
	"context"
	"time"

	"errors"
	"fmt"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/paginator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	shopCollection = "shops"
)

func (repo implRepo) getShopCollection() mongo.Collection {
	return *repo.database.Collection(shopCollection)
}

func (repo implRepo) Create(ctx context.Context, sc models.Scope, opt shop.CreateOption) (models.Shop, error) {
	col := repo.getShopCollection()

	s := repo.buildShopModel(ctx, sc, opt)

	_, err := col.InsertOne(ctx, s)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.Create.InsertOne: %v", err)
		return models.Shop{}, err
	}

	return s, nil
}

func (repo implRepo) Get(ctx context.Context, sc models.Scope, opt shop.GetOption) ([]models.Shop, paginator.Paginator, error) {
	col := repo.getShopCollection()

	filter, err := repo.buildShopQuery(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.Get.buildShopQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	cursor, err := col.Find(ctx, filter, options.Find().
		SetSkip(opt.PagQuery.Offset()).
		SetLimit(opt.PagQuery.Limit).
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1},
		}),
	)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.Find: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var shops []models.Shop
	err = cursor.All(ctx, &shops)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.All: %v", err)
		return nil, paginator.Paginator{}, err
	}

	total, err := col.CountDocuments(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.CountDocuments: %v", err)
		return nil, paginator.Paginator{}, err
	}

	return shops, paginator.Paginator{
		Total:       total,
		Count:       int64(len(shops)),
		PerPage:     opt.PagQuery.Limit,
		CurrentPage: opt.PagQuery.Page,
	}, nil

}

func (repo implRepo) Detail(ctx context.Context, sc models.Scope, id string) (models.Shop, error) {
	col := repo.getShopCollection()

	filter, err := repo.buildShopDetailQuery(ctx, sc, id)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.Detail.buildShopDetailQuery: %v", err)
		return models.Shop{}, err
	}

	var s models.Shop

	err = col.FindOne(ctx, filter).Decode(&s)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.Detail.FindOne: %v", err)
		return models.Shop{}, err
	}

	return s, nil
}
func (repo implRepo) FindByid(ctx context.Context, sc models.Scope, id string) (models.Shop, error) {
	var shop models.Shop
	col := repo.getShopCollection()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return shop, errors.New("invalid id format")
	}

	filter := bson.M{
		"_id":        objectId,
		"deleted_at": nil,
	}

	fmt.Printf("filter: %+v\n", filter)

	err = col.FindOne(ctx, filter).Decode(&shop)
	if err == mongo.ErrNoDocuments {
		return shop, errors.New("shop not found")
	} else if err != nil {
		return shop, err
	}

	return shop, nil
}

func (repo implRepo) Delete(ctx context.Context, sc models.Scope) error {

	col := repo.getShopCollection()
	filter, err := repo.buildShopDetailQuery(ctx, sc, "")
	_, err = col.DeleteOne(ctx, filter)

	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.Detail.DeleteOne: %v", err)
		return err
	}

	return nil
}
func (repo implRepo) Update(ctx context.Context, sc models.Scope, option shop.UpdateOption) (models.Shop, error) {
	col := repo.getShopCollection()
	filter, err := repo.buildShopDetailQuery(ctx, sc, "")
	if err != nil {
		repo.l.Errorf(ctx, "shop.repo.Update.buildshopdetailquery,", err)
		return models.Shop{}, err
	}

	updateData := bson.M{}
	if option.Name != "" {
		updateData["name"] = option.Name
		option.Model.Name = option.Name
	}
	if option.Alias != "" {
		updateData["alias"] = option.Alias
		option.Model.Alias = option.Alias
	}
	if option.City != "" {
		updateData["city"] = option.City
		option.Model.City = option.City
	}
	if option.Street != "" {
		updateData["street"] = option.Street
		option.Model.Street = option.Street
	}
	if option.District != "" {
		updateData["district"] = option.District
		option.Model.District = option.District
	}
	if option.Phone != "" {
		updateData["phone"] = option.Phone
		option.Model.Phone = option.Phone
	}

	updateData["updated_at"] = time.Now()

	update := bson.M{}
	if len(updateData) > 0 {
		update["$set"] = updateData
	}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repo.Update.FindOneAndUpdate:", err)
		return models.Shop{}, err
	}

	return option.Model, nil
}
