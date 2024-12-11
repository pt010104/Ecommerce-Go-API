package redis

import (
	"context"
	"encoding/json"

	"github.com/pt010104/api-golang/internal/models"
)

const (
	userKeyPrefix = "user"
	userExpTime   = 3600
)

func (repo *implRedis) StoreUser(ctx context.Context, user models.User) error {
	userKey := repo.buildUserKey(user.ID.Hex())

	userData, err := json.Marshal(user)
	if err != nil {
		repo.l.Errorf(ctx, "redis.StoreUser.Marshal: %v", err)
		return err
	}

	err = repo.redis.Set(ctx, userKey, userData, userExpTime)
	if err != nil {
		repo.l.Errorf(ctx, "redis.StoreUser.Set: %v", err)
		return err
	}

	return nil
}

func (repo *implRedis) DetailUser(ctx context.Context, userID string) (models.User, error) {
	userKey := repo.buildUserKey(userID)

	data, err := repo.redis.Get(ctx, userKey)
	if err != nil {
		repo.l.Warnf(ctx, "redis.GetUser.Get: %v", err)
		return models.User{}, err
	}

	var user models.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		repo.l.Errorf(ctx, "redis.GetUser.Unmarshal: %v", err)
		return models.User{}, err
	}

	return user, nil
}
