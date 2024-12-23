package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildInventoryModel(context context.Context, opt shop.CreateInventoryOption) (models.Inventory, error) {
	now := util.Now()

	i := models.Inventory{
		ID:         primitive.NewObjectID(),
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

func (impl implRepo) buildInventoryUpdateModel(context context.Context, opt shop.UpdateInventoryOption) (models.Inventory, bson.M, error) {
	now := util.Now()

	setUpdate := bson.M{
		"updated_at": now,
	}
	opt.Model.UpdatedAt = now

	if opt.ReorderLevel != nil && opt.ReorderQuantity != nil {
		setUpdate["reorder_level"] = opt.ReorderLevel
		setUpdate["reorder_quantity"] = opt.ReorderQuantity

		opt.Model.ReorderLevel = opt.ReorderLevel
		opt.Model.ReorderQuantity = opt.ReorderQuantity
	}

	if opt.StockLevel != nil {
		setUpdate["stock_level"] = opt.StockLevel
		opt.Model.StockLevel = *opt.StockLevel
	}

	if opt.ReservedLevel != 0 {
		setUpdate["reserved_level"] = opt.ReservedLevel
		opt.Model.ReservedLevel = opt.ReservedLevel
	}

	var update bson.M
	if len(setUpdate) > 0 {
		update = bson.M{
			"$set": setUpdate,
		}
	}

	return opt.Model, update, nil
}
