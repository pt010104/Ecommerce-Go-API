package http

import (
	"github.com/pt010104/api-golang/internal/cart"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type CartItemReq struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  *int   `json:"quantity" binding:"required"`
}

type UpdateCartRequest struct {
	Item []CartItemReq `json:"items" binding:"required"`
}

func (r UpdateCartRequest) validate() error {
	for _, item := range r.Item {
		if !mongo.IsObjectID(item.ProductID) {
			return errWrongBody
		}
		if item.Quantity == nil || *item.Quantity < 0 {
			return errWrongBody
		}
	}
	return nil
}

func (r UpdateCartRequest) toInput() cart.UpdateInput {
	var items []cart.CartItemInput
	for _, item := range r.Item {
		items = append(items, cart.CartItemInput{
			ProductID: item.ProductID,
			Quantity:  *item.Quantity,
		})
	}
	return cart.UpdateInput{
		NewItemList: items,
	}

}

type shopResponse struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	AvtURL   string            `json:"avt_url,omitempty"`
	Products []productResponse `json:"products"`
}

type productResponse struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type updateCartResponse struct {
	Shops []shopResponse `json:"shops,omitempty"`
}

func (h handler) updateResponse(o cart.UpdateOutput) updateCartResponse {
	shopMap := make(map[string]*shopResponse)
	for _, shop := range o.Shops {
		shopMap[shop.ID.Hex()] = &shopResponse{
			ID:       shop.ID.Hex(),
			Name:     shop.Name,
			Products: []productResponse{},
		}
	}

	for _, cart := range o.Carts {
		shopID := cart.ShopID.Hex()
		if shop, ok := shopMap[shopID]; ok {
			for _, item := range cart.Items {
				shop.Products = append(shop.Products, productResponse{
					ID:       item.ProductID.Hex(),
					Quantity: item.Quantity,
				})
			}
		}
	}

	var response updateCartResponse
	for _, shop := range shopMap {
		response.Shops = append(response.Shops, shopResponse{
			ID:       shop.ID,
			Name:     shop.Name,
			Products: shop.Products,
		})
	}

	return response
}

type ListCartRequest struct {
	UserID  string
	IDs     []string `form:"ids"`
	ShopIDs []string `form:"shop_ids"`
}

// func (r ListCartRequest) toInput() cart.List {
// 	return cart.ListInput{
// 		UserID:  r.UserID,
// 		IDs:     r.IDs,
// 		ShopIDs: r.ShopIDs,
// 	}
// }

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

type DetailCartRequest struct {
	ID     string `uri:"id"`
	UserID string
}

func (r DetailCartRequest) validate() error {
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

type addToCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required"`
}

func (r addToCartRequest) validate() error {
	if !mongo.IsObjectID(r.ProductID) {
		return errWrongBody
	}
	if r.Quantity <= 0 {
		return errWrongBody
	}
	return nil
}

func (r addToCartRequest) toInput() cart.CreateCartInput {
	return cart.CreateCartInput{
		ProductID: r.ProductID,
		Quantity:  r.Quantity,
	}
}

type GetCartReq struct {
	IDs     []string `form:"ids"`
	ShopIDs []string `form:"shop_ids"`
}

func (r GetCartReq) validate() error {
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
func (r GetCartReq) toInput() cart.GetOption {
	return cart.GetOption{
		CartFilter: cart.CartFilter{
			IDs:     r.IDs,
			ShopIDs: r.ShopIDs,
		},
	}
}

type Media_Obj struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}
type GetCartItemResponse struct {
	ProductID string      `json:"product_id"`
	Quantity  int         `json:"quantity"`
	MediaS    []Media_Obj `json:"medias"`
}
type getCartResponseItem struct {
	ID     string                `json:"id"`
	UserID string                `json:"user_id"`
	ShopID string                `json:"shop_id"`
	Item   []GetCartItemResponse `json:"item"`
}
type getCartResponse struct {
	Item []getCartResponseItem `json:"item"`
	Meta listMetaResponse      `json:"meta"`
}
type listMetaResponse struct {
	paginator.PaginatorResponse
}

func (h handler) newGetResponse(carts cart.GetCartOutput) getCartResponse {

	var res getCartResponse
	var resItem []getCartResponseItem
	for _, cart := range carts.CartOutPut {
		var items []GetCartItemResponse
		for _, item := range cart.Cart.Items {
			items = append(items, GetCartItemResponse{
				ProductID: item.ProductID.Hex(),
				Quantity:  item.Quantity,
			})
		}
		resItem = append(resItem, getCartResponseItem{
			ID:     cart.Cart.ID.Hex(),
			UserID: cart.Cart.UserID.Hex(),
			ShopID: cart.Cart.ShopID.Hex(),
			Item:   items,
		})
	}
	res.Item = resItem
	res.Meta = listMetaResponse{
		PaginatorResponse: carts.Paginator.ToResponse(),
	}
	return res

}
