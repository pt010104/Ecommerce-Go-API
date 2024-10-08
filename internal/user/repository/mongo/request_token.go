package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	requestTokenCollection = "request_tokens"
)

func (repo implRepo) getRequestTokenCollection() mongo.Collection {
	return *repo.database.Collection(requestTokenCollection)
}

func (repo implRepo) CreateRequestToken(ctx context.Context, id primitive.ObjectID, token string) (models.RequestToken, error) {
	col := repo.getRequestTokenCollection()
	u, err := repo.buildRequestTokenModel(ctx, id, token)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.buldrequestTokenModel: %v", err)
		return models.RequestToken{}, err

	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.InsertOne: %v", err)
		return models.RequestToken{}, err
	}
	return u, nil

}
