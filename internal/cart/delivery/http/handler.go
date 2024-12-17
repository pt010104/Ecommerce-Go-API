package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/paginator"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary Update a cart
// @Description Update an existing cart with new product details , pass all the item in the cart , once the update api is called , the cart will be replaced with the new items
// @Tags Cart
// @Accept json
// @Produce json
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param request body UpdateCartRequest true "Cart update request"
// @Success 200 {object} updateCartResponse
// @Router /carts/update-cart [post]
func (h handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processUpdateCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Update: %v", err)
		response.Error(c, err)
		return
	}

	o, err := h.uc.Update(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Update: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.updateResponse(o))
}

// @Summary Add a cart
// @Description ADD a product to user cart . Auto group by shop ID
// @Tags Cart
// @Accept json
// @Produce json
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param request body addToCartRequest  true "Add to cart request"
// @Success 200
// @Router /carts/add [post]
func (h handler) Add(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processaddToCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Add: %v", err)
		response.Error(c, err)
		return
	}

	err = h.uc.Add(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Add: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

// @Summary		Get cart
// @Schemes		http https
// @Description	Get shop by ShopIDs, IDs,if no query is passed , all cart of current user will be returned
// @Tags			Cart
// @Accept			json
// @Produce		json
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			ids							query		[]string	false	"IDs"
// @Param			shop_ids				query		[]string		false	"Shop IDs"
// @Param			page						query		int			false	"Page"	default(1)
// @Param			limit						query		int			false	"Limit"	default(10)
// @Success		200							{object}	getCartResponse
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
// @Router			/api/v1/carts [GET]
func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	var pagQuery paginator.PaginatorQuery
	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Get.ShouldBindQuery: %v", err)
		response.Error(c, errWrongPaginationQuery)
		return
	}

	pagQuery.Adjust()
	sc, req, err := h.processGetCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Get: %v", err)
		response.Error(c, err)
		return
	}

	o, err := h.uc.GetCart(ctx, sc, req.toInput(pagQuery))
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Get: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newGetResponse(o))
}
