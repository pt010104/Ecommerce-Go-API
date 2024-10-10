package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildKeyTokenModel(context context.Context, opt user.CreateKeyTokenOption) (models.KeyToken, error) {
	now := time.Now()

	u := models.KeyToken{
		ID:           primitive.NewObjectID(),
		UserID:       opt.UserID,
		SecretKey:    opt.SecretKey,
		RefreshToken: opt.RefrestToken,
		SessionID:    opt.SessionID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return u, nil
}

func (impl implRepo) buildUpdateKeyTokenModel(context context.Context, opt user.UpdateKeyTokenInput) (primitive.M, error) {
	now := time.Now()

	setFields := bson.M{
		"updated_at": now,
	}

	if opt.RefreshToken != "" {
		setFields["refresh_token"] = opt.RefreshToken
	}

	update := bson.M{}
	if len(setFields) > 0 {
		update["$set"] = setFields
	}

	return update, nil
}
