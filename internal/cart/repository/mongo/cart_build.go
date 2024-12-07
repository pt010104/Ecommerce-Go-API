package mongo

import (
	"context"
	"time"

	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (repo implRepo) buildCartModel(opt cart.CreateCartOption, ctx context.Context) (models.Cart, error) {
	now := time.Now()

	p := models.Cart{
		ID:        primitive.NewObjectID(),
		UserID:    opt.UserID,
		ShopID:    opt.ShopID,
		UpdatedAt: now,
		CreatedAt: now,
	}

	return p, nil
}
