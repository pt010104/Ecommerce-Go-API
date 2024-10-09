package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/pt010104/api-golang/internal/user"
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

func (repo implRepo) DetailRequestToken(ctx context.Context, token string) (models.RequestToken, error) {
	col := repo.getRequestTokenCollection()

	filter := repo.buildRequestTokenDetailQuery(ctx, token)

	var requestToken models.RequestToken

	err := col.FindOne(ctx, filter).Decode(&requestToken)
	if err != nil {
		repo.l.Errorf(ctx, "requestToken.repository.mongo.Detail.col.FindOne: %v", err)
		return models.RequestToken{}, err
	}

	return requestToken, nil
}

func (repo implRepo) UpdateRequestToken(ctx context.Context, opt user.UpdateRequestTokenOption) error {
	col := repo.getRequestTokenCollection()

	filter := repo.buildRequestTokenDetailQuery(ctx, opt.Token)

	update, err := repo.buildUpdateRequestTokenModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "requestToken.repository.mongo.UpdateRequestToken.buildUpdateRequestTokenModel: %v", err)
		return err
	}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		repo.l.Errorf(ctx, "requestToken.repository.mongo.UpdateRequestToken.col.UpdateOne: %v", err)
		return err
	}

	return nil
}
