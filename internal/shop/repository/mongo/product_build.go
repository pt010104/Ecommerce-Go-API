package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildProductModel(opt shop.CreateProductOption, ctx context.Context) (models.Product, error) {
	now := time.Now()

	p := models.Product{
		ID:          primitive.NewObjectID(),
		InventoryID: opt.InventoryID,
		Name:        opt.Name,
		Price:       opt.Price,
		CreatedAt:   now,
		UpdatedAt:   now,
		ShopID:      opt.ShopID,
	}

	return p, nil
}
