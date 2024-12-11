package redis

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
)

const (
	sessionKeyPrefix = "session"
	secretKeyPrefix  = "secretkey"
)

func (r implRedis) GetSecretKey(ctx context.Context, sc models.Scope) (string, error) {
	key := r.buildSecretKeyKey(sc.UserID, sc.SessionID)

	result, err := r.redis.Get(ctx, key)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (r implRedis) StoreSecretKey(sc models.Scope, secretKey string, ctx context.Context) error {

	key := r.buildSecretKeyKey(sc.UserID, sc.SessionID)

	return r.redis.Set(ctx, key, secretKey, 0)
}
