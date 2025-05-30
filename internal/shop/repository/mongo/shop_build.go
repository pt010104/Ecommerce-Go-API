package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildShopModel(ctx context.Context, sc models.Scope, opt shop.CreateShopOption) models.Shop {
	now := util.Now()

	s := models.Shop{
		ID:         primitive.NewObjectID(),
		UserID:     mongo.ObjectIDFromHexOrNil(sc.UserID),
		Name:       opt.Name,
		Alias:      opt.Alias,
		City:       opt.City,
		Street:     opt.Street,
		District:   opt.District,
		Phone:      opt.Phone,
		AvgRate:    0,
		IsVerified: false,
		UpdatedAt:  now,
		CreatedAt:  now,
	}

	return s
}
