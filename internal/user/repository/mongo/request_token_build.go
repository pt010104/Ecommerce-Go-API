package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildRequestTokenModel(context context.Context, userId primitive.ObjectID, token string) (models.RequestToken, error) {
	now := util.Now()

	u := models.RequestToken{
		ID:        primitive.NewObjectID(),
		Token:     token,
		UserID:    userId,
		CreatedAt: now,
		UpdateAt:  now,
		IsUsed:    false,
	}
	return u, nil
}

func (repo implRepo) buildUpdateRequestTokenModel(ctx context.Context, opt user.UpdateRequestTokenOption) (bson.M, error) {
	setFields := bson.M{
		"updated_at": util.Now(),
	}

	if opt.IsUsed != nil {
		setFields["is_used"] = *opt.IsUsed
	}

	update := bson.M{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}

	return update, nil
}
