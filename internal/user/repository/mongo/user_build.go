package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
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
	}

	return u, nil
}

func (impl implRepo) buildUpdateUserModel(context context.Context, opt user.UpdateUserOption) (bson.M, models.User, error) {
	setFields := bson.M{
		"updated_at": time.Now(),
	}

	if opt.Email != "" {
		setFields["email"] = opt.Email
		opt.Model.Email = opt.Email
	}
	if opt.IsVerified {
		setFields["is_verified"] = opt.IsVerified
		opt.Model.IsVerified = opt.IsVerified
	}
	if opt.Name != "" {
		setFields["name"] = opt.Name
		opt.Model.Name = opt.Name
	}

	if opt.Password != "" {
		setFields["password"] = opt.Password
		opt.Model.Password = opt.Password
	}

	update := bson.M{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}

	return update, opt.Model, nil
}
