package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	categorycollection = "categories"
)

func (repo implRepo) getCategoryCollection() mongo.Collection {
	return *repo.database.Collection(categorycollection)
}
func (repo implRepo) CreateCategory(ctx context.Context, sc models.Scope, opt admins.CreateCategoryOption) (models.Category, error) {
	col := repo.getCategoryCollection()
	cate := repo.buildCategortModel(ctx, sc, opt)
	_, err := col.InsertOne(ctx, cate)
	if err != nil {
		repo.l.Errorf(ctx, "admins.repo.CreateCategory.Insertone:", err)
		return models.Category{}, err
	}
	return cate, nil
}

func (repo implRepo) ListCategories(ctx context.Context, sc models.Scope, opt admins.GetCategoriesFilter) ([]models.Category, error) {
	col := repo.getCategoryCollection()

	filter := repo.buildCategoryQuery(opt)

	var cates []models.Category
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "admins.repo.ListCategories.Find:", err)
		return nil, err
	}
	err = cursor.All(ctx, &cates)
	if err != nil {
		repo.l.Errorf(ctx, "admins.repo.ListCategories.All:", err)
		return nil, err
	}
	return cates, nil
}
