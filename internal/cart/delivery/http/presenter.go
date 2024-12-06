package http

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/pkg/mongo"
)

type CreateCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	ShopID    string `json:"shop_id"`
	UserID    string `json:"user_id"`
}

func (r CreateCartRequest) validate() error {

	for _, id := range []string{r.ProductID, r.ShopID, r.UserID} {

		if !mongo.IsObjectID(id) {
			return errWrongBody
		}
	}
	return nil

}
func (r CreateCartRequest) toInput() cart.CreateCartInput {
	return cart.CreateCartInput{
		UserID: r.UserID,
		ShopID: r.ShopID,
		Item: []cart.CreateCartItemInput{
			{
				ProductID: r.ProductID,
				Quantity:  r.Quantity,
			},
		},
	}
}
func (r CreateCartRequest) toItemInput() cart.CreateCartItemInput {
	return cart.CreateCartItemInput{
		ProductID: r.ProductID,
		Quantity:  r.Quantity,
	}
}

type CreateCartItemResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
type CreateCartResponse struct {
	ID     string                   `json:"id"`
	UserID string                   `json:"user_id"`
	ShopID string                   `json:"shop_id"`
	Item   []CreateCartItemResponse `json:"item"`
}

func (r CreateCartResponse) toOutput() cart.CreateCartOutput {
	return cart.CreateCartOutput{

		ID:     r.ID,
		UserID: r.UserID,
		ShopID: r.ShopID,
		Item:   []cart.CreateCartItemOutput{},
	}
}
