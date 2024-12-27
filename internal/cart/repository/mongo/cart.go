package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	cartCollection = "carts"
)

func (repo implRepo) getCartCollection() mongo.Collection {
	return *repo.database.Collection(cartCollection)
}

func (repo implRepo) Create(ctx context.Context, sc models.Scope, opt cart.CreateCartOption) (models.Cart, error) {
	col := repo.getCartCollection()
	newCart, err := repo.buildCartModel(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Create.buildCartModel", err)
		return models.Cart{}, err
	}

	_, err = col.InsertOne(ctx, newCart)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Create.InsertOne", err)
		return models.Cart{}, err
	}

	return newCart, nil
}

func (repo implRepo) GetOne(ctx context.Context, sc models.Scope, opt cart.GetOneOption) (models.Cart, error) {
	col := repo.getCartCollection()
	filter, err := repo.buildCartQuery(ctx, sc, opt.CartFilter)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.GetOne.buildCartQuery", err)
		return models.Cart{}, err
	}

	var cart models.Cart
	err = col.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Cart{}, err
		}
		repo.l.Errorf(ctx, "Cart.Repo.GetOne.FindOne", err)
		return models.Cart{}, err
	}
	return cart, nil
}

func (repo implRepo) Update(ctx context.Context, sc models.Scope, opt cart.UpdateCartOption) (models.Cart, error) {
	col := repo.getCartCollection()
	filter, err := repo.buildCartDetailQuery(ctx, sc, opt.Model.ID)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Update.buildCartDetailQuery", err)
		return models.Cart{}, err
	}

	nm, update, err := repo.buildCartUpdateModel(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "Cart.Repo.Update.buildCartUpdateModel", err)
		return models.Cart{}, err
	}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Cart{}, err
		}
		repo.l.Errorf(ctx, "Cart.Repo.Update.FindOneAndUpdate", err)
		return models.Cart{}, err
	}

	return nm, nil

}
func (repo implRepo) ListCart(ctx context.Context, sc models.Scope, opt cart.ListOption) ([]models.Cart, error) {
	col := repo.getCartCollection()

	filter, err := repo.buildCartQuery(ctx, sc, opt.CartFilter)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.buildCartQuery: %v", err)
		return nil, err
	}

	var carts []models.Cart
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.ListCart.Find: %v", err)
		return nil, err
	}

	err = cursor.All(ctx, &carts)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.ListCart.All: %v", err)
		return nil, err
	}
	return carts, nil
}
func (repo implRepo) GetCart(ctx context.Context, sc models.Scope, opt cart.GetOption) ([]models.Cart, paginator.Paginator, error) {
	col := repo.getCartCollection()

	filter, err := repo.buildCartQuery(ctx, sc, opt.CartFilter)
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.Get.buildCartQuery: %v", err)
		return nil, paginator.Paginator{}, err
	}

	cursor, err := col.Find(ctx, filter, options.Find().
		SetSkip(opt.PagQuery.Offset()).
		SetLimit(opt.PagQuery.Limit).
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1},
		}),
	)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.Find: %v", err)
		return nil, paginator.Paginator{}, err
	}

	var carts []models.Cart
	err = cursor.All(ctx, &carts)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.All: %v", err)
		return nil, paginator.Paginator{}, err
	}

	total, err := col.CountDocuments(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "recruitment.candidate.repository.mongo.GetCandidates.CountDocuments: %v", err)
		return nil, paginator.Paginator{}, err
	}

	return carts, paginator.Paginator{
		Total:       total,
		Count:       int64(len(carts)),
		PerPage:     opt.PagQuery.Limit,
		CurrentPage: opt.PagQuery.Page,
	}, nil

}

func (repo implRepo) Delete(ctx context.Context, sc models.Scope, id primitive.ObjectID) error {
	col := repo.getCartCollection()

	_, err := col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		repo.l.Errorf(ctx, "cart.repository.mongo.Delete.DeleteOne: %v", err)
		return err
	}
	return nil
}
