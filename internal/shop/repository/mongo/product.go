package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/paginator"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (repo implRepo) Detailproduct(ctx context.Context, id primitive.ObjectID) (models.Product, error) {
	col := repo.getProductCollection()

	filter, err := repo.buildProductDetailQuery(ctx, id)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.DetailProduct.buildProductDetailQuery: %v", err)
		return models.Product{}, err
	}
	repo.l.Infof(ctx, "Product filter: %v", filter)
	var u models.Product
	err = col.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.DetailProduct.FindOne: %v", err)
		return models.Product{}, err
	}

	return u, nil
}

func (repo implRepo) ListProduct(ctx context.Context, sc models.Scope, opt shop.GetProductFilter) ([]models.Product, error) {
	col := repo.getProductCollection()

	filter, err := repo.buildProductQuery(sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.buildProductQuery: %v", err)
		return []models.Product{}, err
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListProduct.Find: %v", err)
		return []models.Product{}, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListProduct.All: %v", err)
		return []models.Product{}, err
	}

	return products, nil
}

func (repo implRepo) Delete(ctx context.Context, sc models.Scope, ids []string) (err error) {

	col := repo.getProductCollection()
	filter, err := repo.buildProductDeleteQuery(sc, ctx, ids)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.buildProductDeleteQuery: %v", err)
		return err
	}

	_, err = col.DeleteMany(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.Delete.Deletemany: %v", err)
		return err
	}
	return nil

}
func (repo implRepo) GetProduct(ctx context.Context, sc models.Scope, opt shop.GetProductOption) ([]models.Product, paginator.Paginator, error) {
	col := repo.getProductCollection()

	filter, err := repo.buildProductQuery(sc, opt.GetProductFilter)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.Get.buildProductQuery: %v", err)
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

	var products []models.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.All: %v", err)
		return nil, paginator.Paginator{}, err
	}

	total, err := col.CountDocuments(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.CountDocuments: %v", err)
		return nil, paginator.Paginator{}, err
	}

	return products, paginator.Paginator{
		Total:       total,
		Count:       int64(len(products)),
		PerPage:     opt.PagQuery.Limit,
		CurrentPage: opt.PagQuery.Page,
	}, nil
}
func (repo implRepo) IsValidProductID(ctx context.Context, productID primitive.ObjectID) bool {
	col := repo.getProductCollection()
	filter := bson.M{"_id": productID}
	var product models.Product
	err := col.FindOne(ctx, filter).Decode(&product)
	return err == nil
}
func (repo implRepo) UpdateProduct(ctx context.Context, sc models.Scope, option shop.UpdateProductOption) (models.Product, error) {
	col := repo.getProductCollection()
	filter, err := repo.buildProductDetailQuery(ctx, option.ID)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repo.Update.buildshopdetailquery,", err)
		return models.Product{}, err
	}

	option.Model.ID = option.ID
	updateData := bson.M{}
	if option.Name != "" {
		updateData["name"] = option.Name
		option.Model.Name = option.Name
	}
	if option.Alias != "" {
		updateData["alias"] = option.Alias
		option.Model.Alias = option.Alias
	}

	if len(option.CategoryID) > 0 {
		updateData["categoryid"] = option.CategoryID
		option.Model.CategoryID = option.CategoryID
	}
	if len(option.MediaIDs) > 0 {
		updateData["media_ids"] = option.MediaIDs
		option.Model.MediaIDs = option.MediaIDs
	}
	if option.Price != 0 {
		updateData["price"] = option.Price
		option.Model.Price = option.Price
	}

	updateData["updated_at"] = util.Now()

	update := bson.M{}
	if len(updateData) > 0 {
		update["$set"] = updateData
	}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repo.Update.FindOneAndUpdate:", err)
		return models.Product{}, err
	}

	return option.Model, nil
}
func (repo implRepo) UpdateViewProduct(ctx context.Context, id primitive.ObjectID) error {
	col := repo.getProductCollection()
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"view": 1},
		"$set": bson.M{"updated_at": util.Now()},
	}
	_, err := col.UpdateOne(ctx, filter, update)
	return err
}

func (repo implRepo) GetMostViewedProducts(ctx context.Context, sc models.Scope) ([]models.Product, error) {
	col := repo.getProductCollection()
	filter := bson.M{"view": bson.M{"$gt": 0}}
	cursor, err := col.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "view", Value: -1}, {Key: "_id", Value: -1}}))
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.GetMostViewedProducts.Find: %v", err)
		return []models.Product{}, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.GetMostViewedProducts.All: %v", err)
		return []models.Product{}, err
	}

	return products, nil
}

func (repo implRepo) GetMostSoldProducts(ctx context.Context, sc models.Scope) ([]models.Product, error) {
	col := repo.getProductCollection()
	filter := bson.M{"sold": bson.M{"$gt": 0}}
	cursor, err := col.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "sold", Value: -1}, {Key: "_id", Value: -1}}))
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.GetMostSoldProducts.Find: %v", err)
		return []models.Product{}, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.GetMostSoldProducts.All: %v", err)
		return []models.Product{}, err
	}

	return products, nil
}
