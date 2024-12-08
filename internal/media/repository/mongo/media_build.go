package mongo

import (
	"time"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r implRepository) buildMediaModel(opt media.UploadOption) (models.Media, error) {
	now := time.Now()

	m := models.Media{
		ID:        primitive.NewObjectID(),
		UserID:    opt.UserID,
		ShopID:    opt.ShopID,
		FileName:  opt.FileName,
		Folder:    opt.Folder,
		Status:    models.MediaStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return m, nil
}

func (r implRepository) buildUpdateModel(opt media.UpdateOption) bson.M {
	now := time.Now()
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

	return updateSet
}
