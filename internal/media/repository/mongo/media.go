package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/media"
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collection = "media"
)

func (r implRepository) getCollection(sc models.Scope) mongo.Collection {
	return *r.database.Collection(collection)
}

func (r implRepository) Create(ctx context.Context, sc models.Scope, opt media.UploadOption) error {
	col := r.getCollection(sc)

	m, err := r.buildMediaModel(ctx, opt)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Create.buildMediaModel: %v", err)
		return err
	}

	_, err = col.InsertOne(ctx, m)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Create.InsertOne: %v", err)
		return err
	}

	return nil
}

func (r implRepository) Update(ctx context.Context, sc models.Scope, id string, opt media.UpdateOption) (models.Media, error) {
	col := r.getCollection(sc)

	filter, err := r.buildDetailQuery(ctx, sc, id)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Update.buildDetailQuery: %v", err)
		return models.Media{}, err
	}

	update := r.buildUpdateModel(opt)

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Update.UpdateOne: %v", err)
		return models.Media{}, err
	}

	return opt.Model, nil
}

func (r implRepository) Detail(ctx context.Context, sc models.Scope, id string) (models.Media, error) {
	col := r.getCollection(sc)

	filter, err := r.buildDetailQuery(ctx, sc, id)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Detail.buildDetailQuery: %v", err)
		return models.Media{}, err
	}

	var m models.Media
	err = col.FindOne(ctx, filter).Decode(&m)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Detail.FindOne: %v", err)
		return models.Media{}, err
	}

	return m, nil
}
