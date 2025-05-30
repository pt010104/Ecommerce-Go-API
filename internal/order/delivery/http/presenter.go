package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/order"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/response"
)

type CreateCheckoutRequest struct {
	ProductIDs []string `json:"product_ids" binding:"required"`
}

func (r CreateCheckoutRequest) validate() error {
	if len(r.ProductIDs) == 0 {
		return errWrongBody
	}

	for _, id := range r.ProductIDs {
		if !mongo.IsObjectID(id) {
			return errWrongBody
		}
	}

	return nil
}

type shopObject struct {
	ShopID   string `json:"shop_id"`
	ShopName string `json:"shop_name"`
}

type productObject struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image,omitempty"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
}

type itemResponse struct {
	ShopObjects shopObject      `json:"shop_objects"`
	ProductList []productObject `json:"product_list"`
	Price       float64         `json:"price"`
}

type checkoutResponse struct {
	CheckoutID string            `json:"checkout_id"`
	Items      []itemResponse    `json:"items"`
	TotalPrice float64           `json:"total_price"`
	ExpiredAt  response.DateTime `json:"expired_at"`
}

func (h handler) newCreateCheckoutCheckoutResponse(o order.CreateCheckoutOutput) checkoutResponse {
	var checkoutResponse checkoutResponse
	image_urls := make(map[string]string)
	for i, url := range o.ImageURLs {
		image_urls[o.Products[i].ID.Hex()] = url
	}

	for _, shop := range o.Shops {
		shopObject := shopObject{
			ShopID:   shop.ID.Hex(),
			ShopName: shop.Name,
		}

		productList := make([]productObject, 0)
		for _, product := range o.Products {
			if product.ShopID.Hex() == shop.ID.Hex() {
				productList = append(productList, productObject{
					ProductID:    product.ID.Hex(),
					ProductName:  product.Name,
					ProductImage: image_urls[product.ID.Hex()],
					Price:        product.Price,
					Quantity:     o.QuantityMap[product.ID.Hex()],
				})
			}
		}

		checkoutResponse.Items = append(checkoutResponse.Items, itemResponse{
			ShopObjects: shopObject,
			ProductList: productList,
			Price:       float64(o.TotalPriceByShop[shop.ID.Hex()]),
		})
	}

	checkoutResponse.TotalPrice = float64(o.TotalPrice)
	checkoutResponse.ExpiredAt = response.DateTime(o.ExpiredAt)
	checkoutResponse.CheckoutID = o.CheckoutID

	return checkoutResponse
}

type CreateOrderRequest struct {
	CheckoutID    string `json:"checkout_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	AddressID     string `json:"address_id" binding:"required"`
	VoucherID     string `json:"voucher_id"`
}

func (r CreateOrderRequest) validate() error {
	if !mongo.IsObjectID(r.CheckoutID) {
		return errWrongBody
	}

	if r.VoucherID != "" && !mongo.IsObjectID(r.VoucherID) {
		return errWrongBody
	}

	return nil
}

func (r CreateOrderRequest) toInput() order.CreateOrderInput {
	return order.CreateOrderInput{
		CheckoutID:    r.CheckoutID,
		PaymentMethod: r.PaymentMethod,
		AddressID:     r.AddressID,
		VoucherID:     r.VoucherID,
	}
}

type createOrderResponse struct {
	OrderID string `json:"order_id"`
}

func (h handler) newCreateOrderResponse(o models.Order) createOrderResponse {
	return createOrderResponse{
		OrderID: o.ID.Hex(),
	}
}

type ListOrderRequest struct {
	Status string `form:"status" binding:"required"`
}

func (r ListOrderRequest) validate() error {
	if r.Status != models.OrderStatusPending && r.Status != models.OrderStatusProcessing && r.Status != models.OrderStatusShipping && r.Status != models.OrderStatusDelivered && r.Status != models.OrderStatusCanceled {
		return errWrongBody
	}

	return nil
}

func (r ListOrderRequest) toInput() order.ListOrderInput {
	return order.ListOrderInput{
		Status: r.Status,
	}
}

type orderResponse struct {
	OrderID    string            `json:"order_id"`
	Status     string            `json:"status"`
	Address    string            `json:"address,omitempty"`
	UserName   string            `json:"user_name,omitempty"`
	UserPhone  string            `json:"user_phone,omitempty"`
	UserEmail  string            `json:"user_email,omitempty"`
	TotalPrice float64           `json:"total_price"`
	Products   []productObject   `json:"products"`
	CreatedAt  response.DateTime `json:"created_at"`
}

func (h handler) newListOrderResponse(o order.ListOrderOutput) []orderResponse {
	resp := make([]orderResponse, 0)
	for _, order := range o.Orders {
		products := make([]productObject, 0)
		for _, p := range order.Products {
			products = append(products, productObject{
				ProductID:   p.ProductID,
				ProductName: p.ProductName,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}

		resp = append(resp, orderResponse{
			OrderID:    order.Order.ID.Hex(),
			Status:     order.Order.Status,
			TotalPrice: order.TotalPrice,
			Products:   products,
			CreatedAt:  response.DateTime(order.Order.CreatedAt),
		})
	}
	return resp
}

type ListOrderShopRequest struct {
	Status string `form:"status" binding:"required"`
}

func (r ListOrderShopRequest) validate() error {
	if r.Status != models.OrderStatusPending && r.Status != models.OrderStatusProcessing && r.Status != models.OrderStatusShipping && r.Status != models.OrderStatusDelivered && r.Status != models.OrderStatusCanceled {
		return errWrongBody
	}

	return nil
}

func (r ListOrderShopRequest) toInput() order.ListOrderShopInput {
	return order.ListOrderShopInput{
		Status: r.Status,
	}
}
func (h handler) newListOrderShopResponse(o order.ListOrderShopOutput) []orderResponse {
	resp := make([]orderResponse, 0)
	for i, order := range o.Orders {
		products := make([]productObject, 0)
		for _, p := range order.Products {
			products = append(products, productObject{
				ProductID:   p.ProductID,
				ProductName: p.ProductName,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}

		resp = append(resp, orderResponse{
			OrderID:    order.Order.ID.Hex(),
			Status:     order.Order.Status,
			TotalPrice: order.TotalPrice,
			Products:   products,
			CreatedAt:  response.DateTime(order.Order.CreatedAt),
			Address:    o.AddressMerges[i],
			UserName:   o.UserModels[i].Name,
			UserPhone:  o.UserModels[i].Address[0].Phone,
			UserEmail:  o.UserModels[i].Email,
		})
	}
	return resp
}

type UpdateOrderRequest struct {
	Status string `json:"status" binding:"required"`
}

func (r UpdateOrderRequest) validate() error {
	if r.Status != models.OrderStatusPending && r.Status != models.OrderStatusProcessing && r.Status != models.OrderStatusShipping && r.Status != models.OrderStatusDelivered && r.Status != models.OrderStatusCanceled {
		return errWrongBody
	}

	return nil
}

func (r UpdateOrderRequest) toInput(orderID string) order.UpdateOrderInput {
	return order.UpdateOrderInput{
		OrderID: orderID,
		Status:  r.Status,
	}
}
