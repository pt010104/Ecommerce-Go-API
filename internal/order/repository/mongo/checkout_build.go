package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCheckoutModel(ctx context.Context, sc models.Scope, opt order.CreateCheckoutOption) (models.Checkout, error) {
	now := util.Now()

	p := models.Checkout{
		ID:        primitive.NewObjectID(),
		Products:  opt.Products,
		UserID:    mongo.ObjectIDFromHexOrNil(sc.UserID),
		Status:    models.CheckoutStatusPending,
		ExpiredAt: now.Add(time.Minute * 10),
		UpdatedAt: now,
		CreatedAt: now,
	}

	return p, nil
}

func (repo implRepo) buildCheckoutUpdate(ctx context.Context, opt order.UpdateCheckoutOption) (models.Checkout, bson.M) {
	now := util.Now()

	setFields := bson.M{
		"updated_at": now,
	}

	if opt.Status != "" {
		setFields["status"] = opt.Status
	}

	opt.Model.UpdatedAt = now
	if opt.Status != "" {
		opt.Model.Status = opt.Status
	}

	update := bson.M{
		"$set": setFields,
	}

	return opt.Model, update
}
