package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl implRepo) buildKeyTokenModel(context context.Context, opt user.CreateKeyTokenOption) (models.KeyToken, error) {
	now := util.Now()

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
	now := util.Now()

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
