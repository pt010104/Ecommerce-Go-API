package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary Create a cart
// @Description Create a new cart with product details
// @Tags Cart
// @Accept json
// @Produce json
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param request body CreateCartRequest true "Cart creation request"
// @Success 200 {object} CreateCartResponse
// @Router /carts [post]
func (h handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Create: %v", err)
		response.Error(c, err)
		return
	}

	cart, err := h.uc.Create(sc, ctx, req.toInput(), req.toItemInput())
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Create: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.CreateResponse(cart))
}

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
// @Success 200 {object} UpdateCartResponse
// @Router /carts/update-cart [post]
func (h handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	_, req, err := h.processUpdateCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Update: %v", err)
		response.Error(c, err)
		return
	}

	cart, err := h.uc.Update(ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Update: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.UpdateResponse(cart))
}

// @Summary List carts
// @Description Retrieve a list of carts based on user ID and filters
// @Tags Cart
// @Accept json
// @Produce json
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param user_id query string true "User ID"
// @Param ids query []string false "Cart IDs (optional)"
// @Param shop_ids query []string false "Shop IDs (optional)"
// @Success 200 {array} ListCartResponse
// @Router /carts [get]
func (h handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processListCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.List: %v", err)
		response.Error(c, err)
		return
	}

	cart, err := h.uc.ListCart(sc, ctx, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.List: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newListResponse(cart))
}
