package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	keyTokenCollection = "key_tokens"
)

func (repo implRepo) getKeyTokenCollection() mongo.Collection {
	return *repo.database.Collection(keyTokenCollection)
}

func (repo implRepo) CreateKeyToken(ctx context.Context, opt user.CreateKeyTokenOption) (models.KeyToken, error) {
	col := repo.getKeyTokenCollection()
	u, err := repo.buildKeyTokenModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.buldUserModel: %v", err)
		return models.KeyToken{}, err

	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.InsertOne: %v", err)
		return models.KeyToken{}, err
	}
	return u, nil

}

func (repo implRepo) DetailKeyToken(ctx context.Context, userID string, sessionID string) (models.KeyToken, error) {
	col := repo.getKeyTokenCollection()

	filter, err := repo.buildKeyTokenDetailQuery(ctx, userID, sessionID)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.DetailKeyToken.buildKeyTokenDetailQuery: %v", err)
		return models.KeyToken{}, err
	}

	var keyToken models.KeyToken

	err = col.FindOne(ctx, filter).Decode(&keyToken)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.DetailKeyToken.col.FindOne: %v", err)
		return models.KeyToken{}, err
	}

	return keyToken, nil
}

func (repo implRepo) DeleteKeyToken(ctx context.Context, userID string, sessionID string) error {
	col := repo.getKeyTokenCollection()

	filter, err := repo.buildKeyTokenDetailQuery(ctx, userID, sessionID)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.DeleteKeyToken.buildKeyTokenDetailQuery: %v", err)
		return err
	}

	_, err = col.DeleteOne(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.DeleteKeyToken.DeleteOne: %v", err)
		return err
	}
	return nil

}
