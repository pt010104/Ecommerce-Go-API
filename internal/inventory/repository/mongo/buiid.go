package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/inventory"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildInventoryModel(context context.Context, opt inventory.CreateOption) (models.Inventory, error) {
	now := time.Now()

	i := models.Inventory{
		ID:         primitive.NewObjectID(),
		ProductID:  mongo.ObjectIDFromHexOrNil(opt.ProductID),
		StockLevel: opt.StockLevel,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if opt.ReorderLevel != nil && opt.ReorderQuantity != nil {
		i.ReorderLevel = opt.ReorderLevel
		i.ReorderQuantity = opt.ReorderQuantity
	}

	return i, nil
}
