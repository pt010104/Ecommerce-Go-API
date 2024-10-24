package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/admin"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCategortModel(ctx context.Context, sc models.Scope, opt admin.CreateCategoryOption) models.Category {

	return models.Category{
		ID:          primitive.NewObjectID(),
		Name:        opt.Name,
		Description: opt.Description,
		CreatedBy:   mongo.ObjectIDFromHexOrNil(sc.UserID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
