package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "users"
)

func (repo implRepo) getUserCollection() mongo.Collection {
	return *repo.database.Collection(userCollection)
}

func (repo implRepo) CreateUserRepo(ctx context.Context, opt user.RepoOption) (models.User, error) {
	col := repo.getUserCollection()

	u, err := repo.buildUserModel(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.buldUserModel: %v", err)
		return models.User{}, err

	}

	_, err = col.InsertOne(ctx, u)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.Create.InsertOne: %v", err)
		return models.User{}, err
	}
	return u, nil
}
