package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildProductModel(opt shop.CreateProductOption, ctx context.Context) (models.Product, error) {
	now := util.Now()

	p := models.Product{
		ID:          primitive.NewObjectID(),
		InventoryID: opt.InventoryID,
		Name:        opt.Name,
		Price:       opt.Price,
		CreatedAt:   now,
		UpdatedAt:   now,
		ShopID:      opt.ShopID,
		CategoryID:  opt.CategoryID,
		Alias:       opt.Alias,
		MediaIDs:    opt.MediaIDs,
	}

	return p, nil
}
