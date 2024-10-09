package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
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
