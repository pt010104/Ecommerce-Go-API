package http

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
)

type CreateCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	ShopID    string `json:"shop_id"`
	UserID    string
}

func (r CreateCartRequest) validate() error {

	for _, id := range []string{r.ProductID, r.ShopID} {

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
func (h handler) CreateResponse(c models.Cart) CreateCartResponse {
	var items []CreateCartItemResponse
	for _, item := range c.Items {
		items = append(items, CreateCartItemResponse{
			ProductID: item.ProductID.Hex(),
			Quantity:  item.Quantity,
		})
	}
	return CreateCartResponse{
		ID:     c.ID.Hex(),
		UserID: c.UserID.Hex(),
		ShopID: c.ShopID.Hex(),
		Item:   items,
	}
}

type UpdateCartItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
type UpdateCartRequest struct {
	ID     string `json:"id"`
	UserID string
	ShopID string                  `json:"shop_id"`
	Item   []UpdateCartItemRequest `json:"item"`
}

func (r UpdateCartRequest) validate() error {
	if !mongo.IsObjectID(r.ID) {
		return errWrongBody
	}
	if !mongo.IsObjectID(r.ShopID) {
		return errWrongBody
	}
	for _, item := range r.Item {
		if !mongo.IsObjectID(item.ProductID) {
			return errWrongBody
		}
	}
	return nil
}

type UpdateCartItemResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
type UpdateCartResponse struct {
	ID     string                   `json:"id"`
	UserID string                   `json:"user_id"`
	ShopID string                   `json:"shop_id"`
	Item   []UpdateCartItemResponse `json:"item"`
}

func (r UpdateCartResponse) toOutput() UpdateCartResponse {
	return UpdateCartResponse{
		ID:     r.ID,
		UserID: r.UserID,
		ShopID: r.ShopID,
		Item:   []UpdateCartItemResponse{},
	}
}
func (h handler) UpdateResponse(c models.Cart) UpdateCartResponse {

	var items []UpdateCartItemResponse
	for _, item := range c.Items {
		items = append(items, UpdateCartItemResponse{
			ProductID: item.ProductID.Hex(),
			Quantity:  item.Quantity,
		})
	}
	return UpdateCartResponse{
		ID:     c.ID.Hex(),
		UserID: c.UserID.Hex(),
		ShopID: c.ShopID.Hex(),
		Item:   items,
	}
}
func (r UpdateCartRequest) toInput() cart.UpdateCartOption {
	var items []models.CartItem
	for _, item := range r.Item {
		items = append(items, models.CartItem{
			ProductID: mongo.ObjectIDFromHexOrNil(item.ProductID),
			Quantity:  item.Quantity,
		})
	}
	return cart.UpdateCartOption{
		ID:          mongo.ObjectIDFromHexOrNil(r.ID),
		UserID:      mongo.ObjectIDFromHexOrNil(r.UserID),
		ShopID:      mongo.ObjectIDFromHexOrNil(r.ShopID),
		NewItemList: items,
	}

}

type ListCartRequest struct {
	UserID  string
	IDs     []string `form:"ids"`
	ShopIDs []string `form:"shop_ids"`
}

func (r ListCartRequest) toInput() cart.GetCartFilter {
	return cart.GetCartFilter{
		UserID:  r.UserID,
		IDs:     r.IDs,
		ShopIDs: r.ShopIDs,
	}
}
func (r ListCartRequest) validate() error {
	if !mongo.IsObjectID(r.UserID) {
		return errWrongBody
	}
	if len(r.IDs) > 0 {
		for _, id := range r.IDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}
	if len(r.ShopIDs) > 0 {
		for _, id := range r.ShopIDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}
	return nil

}

type ListCartResponse struct {
	ID     string                 `json:"id"`
	UserID string                 `json:"user_id"`
	ShopID string                 `json:"shop_id"`
	Item   []ListCartItemResponse `json:"item"`
}
type ListCartItemResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type GetCartRequest struct {
	ID     string `uri:"id"`
	UserID string
}

func (r GetCartRequest) validate() error {
	if !mongo.IsObjectID(r.ID) {
		return errWrongBody
	}
	return nil
}
func (h handler) newListResponse(carts []models.Cart) []ListCartResponse {
	var res []ListCartResponse
	for _, cart := range carts {
		var items []ListCartItemResponse
		for _, item := range cart.Items {
			items = append(items, ListCartItemResponse{
				ProductID: item.ProductID.Hex(),
				Quantity:  item.Quantity,
			})
		}
		res = append(res, ListCartResponse{
			ID:     cart.ID.Hex(),
			UserID: cart.UserID.Hex(),
			ShopID: cart.ShopID.Hex(),
			Item:   items,
		})
	}
	return res
}
