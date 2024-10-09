package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildRequestTokenModel(context context.Context, userId primitive.ObjectID, token string) (models.RequestToken, error) {
	now := time.Now()

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
		"updated_at": time.Now(),
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
