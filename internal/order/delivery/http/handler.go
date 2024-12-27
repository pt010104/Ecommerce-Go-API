package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary		Create checkout
// @Schemes		http https
// @Description	Create checkout
// @Tags			Order
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			request body CreateCheckoutRequest true "Request Body"
//
// @Success		200							{object}	checkoutResponse	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/orders/checkout [POST]
func (h *handler) CreateCheckout(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateCheckoutRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.CreateCheckout: %v", err)
		response.Error(c, err)
		return
	}

	o, err := h.uc.CreateCheckout(ctx, sc, req.ProductIDs)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.CreateCheckout: %v", err)
		response.Error(c, h.mapErrors(err))
		return
	}

	response.OK(c, h.newCreateCheckoutCheckoutResponse(o))
}

// @Summary		Create order
// @Schemes		http https
// @Description	Create order
// @Tags			Order
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			request body CreateOrderRequest true "Request Body"
//
// @Success		200							{object}	response.Resp	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/orders [POST]
func (h *handler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateOrderRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.CreateOrder: %v", err)
		response.Error(c, err)
		return
	}

	o, err := h.uc.CreateOrder(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.CreateOrder: %v", err)
		response.Error(c, h.mapErrors(err))
		return
	}

	response.OK(c, h.newCreateOrderResponse(o))
}

// @Summary		List order
// @Schemes		http https
// @Description	List order
// @Tags			Order
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			request query ListOrderRequest true "Request Body"
//
// @Success		200							{object}	response.Resp	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/orders [GET]
func (h *handler) ListOrder(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processListOrderRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.ListOrder: %v", err)
		response.Error(c, err)
		return
	}

	orders, err := h.uc.ListOrder(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.ListOrder: %v", err)
		response.Error(c, h.mapErrors(err))
		return
	}

	response.OK(c, h.newListOrderResponse(orders))
}

// @Summary		List order shop
// @Schemes		http https
// @Description	List orders for shop
// @Tags			Order
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			request query ListOrderShopRequest true "Request Body"
//
// @Success		200							{object}	response.Resp	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/orders/shop [GET]
func (h *handler) ListOrderShop(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processListOrderShopRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.ListOrderShop: %v", err)
		response.Error(c, err)
		return
	}

	orders, err := h.uc.ListOrderShop(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.ListOrderShop: %v", err)
		response.Error(c, h.mapErrors(err))
		return
	}

	response.OK(c, h.newListOrderShopResponse(orders))
}

// @Summary		Update order
// @Schemes		http https
// @Description	Update order status
// @Tags			Order
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			order_id					path		string		true	"Order ID"
// @Param			request body UpdateOrderRequest true "Request Body"
//
// @Success		200							{object}	response.Resp	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/orders/shop/{order_id} [PATCH]
func (h *handler) UpdateOrder(c *gin.Context) {
	ctx := c.Request.Context()

	sc, orderID, req, err := h.processUpdateOrderRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.UpdateOrder: %v", err)
		response.Error(c, err)
		return
	}

	err = h.uc.UpdateOrder(ctx, sc, req.toInput(orderID))
	if err != nil {
		h.l.Errorf(ctx, "order.delivery.http.handler.UpdateOrder: %v", err)
		response.Error(c, h.mapErrors(err))
		return
	}

	response.OK(c, nil)
}
