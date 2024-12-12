package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildUserModel(context context.Context, opt user.CreateUserOption) (models.User, error) {
	now := time.Now()

	u := models.User{
		ID:         primitive.NewObjectID(),
		Email:      opt.Email,
		Name:       opt.Name,
		Password:   opt.Password,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsVerified: false,
		Role:       0,
		MediaID:    primitive.NilObjectID,
	}

	return u, nil
}

func (impl implRepo) buildUpdateUserModel(context context.Context, opt user.UpdateUserOption) (bson.M, models.User, error) {
	setFields := bson.M{
		"name":       opt.Name,
		"email":      opt.Email,
		"updated_at": time.Now(),
	}
	opt.Model.Name = opt.Name
	opt.Model.Email = opt.Email
	opt.Model.UpdatedAt = time.Now()

	if opt.IsVerified {
		setFields["is_verified"] = opt.IsVerified
		opt.Model.IsVerified = opt.IsVerified
	}

	if opt.Password != "" {
		setFields["password"] = opt.Password
		opt.Model.Password = opt.Password
	}

	if opt.MediaID != "" {
		setFields["media_id"] = opt.MediaID
		opt.Model.MediaID = mongo.ObjectIDFromHexOrNil(opt.MediaID)
	}

	update := bson.M{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}

	return update, opt.Model, nil
}
