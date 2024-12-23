package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCategortModel(ctx context.Context, sc models.Scope, opt admins.CreateCategoryOption) models.Category {

	return models.Category{
		ID:          primitive.NewObjectID(),
		Name:        opt.Name,
		Description: opt.Description,
		CreatedBy:   mongo.ObjectIDFromHexOrNil(sc.UserID),
		CreatedAt:   util.Now(),
		UpdatedAt:   util.Now(),
	}
}
