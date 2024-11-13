package mongo

import (
	"context"
	"fmt"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	productCollection = "products"
)

func (repo implRepo) getProductCollection() mongo.Collection {
	return *repo.database.Collection(productCollection)
}

const (
	CategoryCollection = "categories"
)

func (repo implRepo) getCategoryCollection() mongo.Collection {
	return *repo.database.Collection(CategoryCollection)
}
func (repo implRepo) ValidateCategoryIDs(ctx context.Context, categoryIDs []primitive.ObjectID) error {
	colC := repo.getCategoryCollection()

	filter := bson.M{"_id": bson.M{"$in": categoryIDs}}
	count, err := colC.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count != int64(len(categoryIDs)) {
		return fmt.Errorf("some category IDs are invalid")
	}

	return nil
}
func (repo implRepo) CreateProduct(ctx context.Context, sc models.Scope, opt shop.CreateProductOption) (models.Product, error) {
	colP := repo.getProductCollection()
	err1 := repo.ValidateCategoryIDs(ctx, opt.CategoryID)
	if err1 != nil {
		repo.l.Errorf(ctx, "shop.repo.product.validatecateIDS:", err1)
		return models.Product{}, shop.ErrNonExistCategory
	}
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

	filter, err := repo.buildProductQuery(opt)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListProduct.buildProductQuery: %v", err)
		return nil, err
	}

	fmt.Println(filter)

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListProduct.Find: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	err = cursor.All(ctx, &products)
	if err != nil {
		repo.l.Errorf(ctx, "shop.repository.mongo.ListProduct.All: %v", err)
		return nil, err
	}

	return products, nil
}
