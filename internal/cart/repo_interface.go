package cart

import (
	"context"

	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repo interface {
	Create(opt CreateCartOption, opt2 CreateCartItemOption, ctx context.Context) (models.Cart, error)
	Get(ctx context.Context, ID primitive.ObjectID) (models.Cart, error)
	Update(ctx context.Context, opt UpdateCartOption) (models.Cart, error)
	ListCart(sc models.Scope, ctx context.Context, opt GetCartFilter) ([]models.Cart, error)
}
