package cart

import (
	"github.com/pt010104/api-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCartOption struct {
	UserID primitive.ObjectID
	ShopID primitive.ObjectID
}
type CreateCartItemOption struct {
	ProductID primitive.ObjectID
	Quantity  int
}
type UpdateCartOption struct {
	Model       models.Cart
	ID          primitive.ObjectID
	ShopID      primitive.ObjectID
	UserID      primitive.ObjectID
	NewItemList []models.CartItem
}
type UpdateCartItemOption struct {
	Quantity int
}
type GetCartFilter struct {
	UserID  string
	IDs     []string
	ShopIDs []string
}
