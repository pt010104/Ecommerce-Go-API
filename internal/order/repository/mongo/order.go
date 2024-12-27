package mongo

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	orderCollection = "orders"
)

func (repo implRepo) getOrderCollection() mongo.Collection {
	return *repo.database.Collection(orderCollection)
}

func (repo implRepo) CreateOrder(ctx context.Context, sc models.Scope, opt order.CreateOrderOption) (models.Order, error) {
	col := repo.getOrderCollection()
	newOrder, err := repo.buildOrderModel(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.CreateOrder.buildOrderModel", err)
		return models.Order{}, err
	}

	_, err = col.InsertOne(ctx, newOrder)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.CreateOrder.InsertOne", err)
		return models.Order{}, err
	}

	return newOrder, nil
}

func (repo implRepo) DetailOrder(ctx context.Context, sc models.Scope, orderID string) (models.Order, error) {
	col := repo.getOrderCollection()

	filter, err := repo.buildOrderDetailQuery(ctx, sc, orderID)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.DetailOrder.buildOrderDetailQuery", err)
		return models.Order{}, err
	}

	var order models.Order
	err = col.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.DetailOrder.FindOne", err)
		return models.Order{}, err
	}

	return order, nil
}

func (repo implRepo) ListOrder(ctx context.Context, sc models.Scope, opt order.ListOrderOption) ([]models.Order, error) {
	col := repo.getOrderCollection()

	filter, err := repo.buildOrderQuery(ctx, sc, opt)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.ListOrder.buildOrderQuery", err)
		return []models.Order{}, err
	}

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.ListOrder.Find", err)
		return []models.Order{}, err
	}

	var orders []models.Order
	err = cursor.All(ctx, &orders)
	if err != nil {
		repo.l.Errorf(ctx, "Order.Repo.ListOrder.All", err)
		return []models.Order{}, err
	}

	return orders, nil
}
