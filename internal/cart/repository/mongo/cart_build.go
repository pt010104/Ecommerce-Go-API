package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCartModel(ctx context.Context, sc models.Scope, opt cart.CreateCartOption) (models.Cart, error) {
	now := time.Now()

	for _, item := range opt.CartItemList {
		item.UpdatedAt = now
	}

	p := models.Cart{
		ID:        primitive.NewObjectID(),
		UserID:    mongo.ObjectIDFromHexOrNil(sc.UserID),
		ShopID:    mongo.ObjectIDFromHexOrNil(opt.ShopID),
		Items:     opt.CartItemList,
		UpdatedAt: now,
		CreatedAt: now,
	}

	return p, nil
}

func (repo implRepo) buildCartUpdateModel(ctx context.Context, sc models.Scope, opt cart.UpdateCartOption) (models.Cart, bson.M, error) {
	now := time.Now()

	unsetFields := bson.M{}

	setFields := bson.M{
		"updated_at": now,
	}
	opt.Model.UpdatedAt = now

	if len(opt.NewItemList) > 0 {
		setFields["items"] = opt.NewItemList
	} else {
		unsetFields["items"] = ""
	}
	opt.Model.Items = opt.NewItemList

	update := bson.M{
		"$set": setFields,
	}

	if len(unsetFields) > 0 {
		update["$unset"] = unsetFields
	}

	return opt.Model, update, nil
}
