package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/checkout"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCheckoutModel(ctx context.Context, sc models.Scope, opt checkout.CreateOption) (models.Checkout, error) {
	now := util.Now()

	p := models.Checkout{
		ID:        primitive.NewObjectID(),
		CartIDs:   opt.CartIDs,
		UserID:    mongo.ObjectIDFromHexOrNil(sc.UserID),
		Status:    models.CheckoutStatusPending,
		ExpiredAt: now.Add(time.Minute * 10),
		UpdatedAt: now,
		CreatedAt: now,
	}

	return p, nil
}
