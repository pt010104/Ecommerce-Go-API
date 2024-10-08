package mongo

import (
	"context"

	"fmt"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollection = "users"
)

func (repo implRepo) getUserCollection() mongo.Collection {
	return *repo.database.Collection(userCollection)
}

func (repo implRepo) CreateUser(ctx context.Context, opt user.CreateUserOption) (models.User, error) {
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

func (repo implRepo) GetUser(ctx context.Context, opt user.GetUserOption) (models.User, error) {
	col := repo.getUserCollection()

	filter, err := repo.buidUserQuery(ctx, opt)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.GetUser.buidUserQuery: %v", err)
		return models.User{}, err
	}

	var user models.User

	err = col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		repo.l.Errorf(ctx, "survey.repository.mongo.GetUser.col.FindOne: %v", err)
		return models.User{}, err
	}

	return user, nil
}
func (repo implRepo) UpdateRecord(ctx context.Context, UserID string, updateData bson.M) error {

	col := repo.getUserCollection()

	filter := bson.M{"_id": UserID}

	update := bson.M{
		"$set": updateData,
	}

	_, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}

	return nil
}
func (repo implRepo) UpdateRequestTokenRecord(ctx context.Context, JWT string, updateData bson.M) error {

	col := repo.getRequestTokenCollection()

	filter := bson.M{"token": JWT}

	update := bson.M{
		"$set": updateData,
	}

	_, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update record: %v", err)
	}

	return nil
}
func (repo implRepo) DetailUser(ctx context.Context, id string) (models.User, error) {
	col := repo.getUserCollection()

	filter, err := repo.buildUserDetailQuery(ctx, id)
	if err != nil {
		repo.l.Errorf(ctx, "user.repository.mongo.GetUser.buildUserDetailQuery: %v", err)
		return models.User{}, err
	}

	var user models.User

	err = col.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		repo.l.Errorf(ctx, "survey.repository.mongo.Detail.col.FindOne: %v", err)
		return models.User{}, err
	}

	return user, nil
}
func (repo implRepo) IsJWTresetVaLID(ctx context.Context, JWT string) (bool, error) {

	col := repo.getRequestTokenCollection()

	filter := bson.M{"token": JWT}

	var tokenRecord struct {
		Token  string `bson:"token"`
		IsUsed bool   `bson:"is_used"`
	}

	err := col.FindOne(ctx, filter).Decode(&tokenRecord)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return false, nil
		}

		return false, err
	}
	if tokenRecord.IsUsed {
		return false, nil
	}
	return true, nil
}
