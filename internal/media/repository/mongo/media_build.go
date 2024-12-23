package mongo

import (
	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r implRepository) buildMediaModel(sc models.Scope, opt media.UploadOption) (models.Media, error) {
	now := util.Now()

	m := models.Media{
		ID:        primitive.NewObjectID(),
		UserID:    mongo.ObjectIDFromHexOrNil(sc.UserID),
		ShopID:    mongo.ObjectIDFromHexOrNil(sc.ShopID),
		FileName:  opt.FileName,
		Folder:    opt.Folder,
		Status:    models.MediaStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return m, nil
}

func (r implRepository) buildUpdateModel(opt media.UpdateOption) bson.M {
	now := util.Now()
	updateSet := bson.M{}

	if opt.Status != "" {
		updateSet["status"] = opt.Status
		opt.Model.Status = opt.Status
	}

	if opt.URL != "" {
		updateSet["url"] = opt.URL
		opt.Model.URL = opt.URL
	}

	updateSet["updated_at"] = now
	opt.Model.UpdatedAt = now

	return bson.M{
		"$set": updateSet,
	}
}
