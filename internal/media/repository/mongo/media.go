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

func (r implRepository) Create(ctx context.Context, sc models.Scope, opt media.UploadOption) (models.Media, error) {
	col := r.getCollection(sc)

	m, err := r.buildMediaModel(sc, opt)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Create.buildMediaModel: %v", err)
		return models.Media{}, err
	}

	_, err = col.InsertOne(ctx, m)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.Create.InsertOne: %v", err)
		return models.Media{}, err
	}

	return m, nil
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

func (r implRepository) List(ctx context.Context, sc models.Scope, opt media.ListOption) ([]models.Media, error) {
	col := r.getCollection(sc)

	filter, err := r.buildQuery(ctx, sc, opt.GetFilter)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.List.buildQuery: %v", err)
		return nil, err
	}

	var medias []models.Media
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.List.Find: %v", err)
		return nil, err
	}

	err = cursor.All(ctx, &medias)
	if err != nil {
		r.l.Errorf(ctx, "media.repository.mongo.List.All: %v", err)
		return nil, err
	}

	return medias, nil
}
