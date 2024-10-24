package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	categorycollection = "categories"
)

func (repo implRepo) getCategoryCollection() mongo.Collection {
	return *repo.database.Collection(categorycollection)
}
func (repo implRepo) CreateCategory(ctx context.Context, sc models.Scope, opt admin.CreateCategoryOption) (models.Category, error) {
	col := repo.getCategoryCollection()
	cate := repo.buildCategortModel(ctx, sc, opt)
	_, err := col.InsertOne(ctx, cate)
	if err != nil {
		repo.l.Errorf(ctx, "admin.repo.CreateCategory.Insertone:", err)
		return models.Category{}, err
	}
	return cate, nil
}
