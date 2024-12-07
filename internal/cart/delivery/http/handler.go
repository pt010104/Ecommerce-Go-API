package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

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

	response.OK(c, cart)
}
func (h handler) Get(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processGetCartRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Get: %v", err)
		response.Error(c, err)
		return
	}

	cart, err := h.uc.GetCart(sc, ctx, req.ID)
	if err != nil {
		h.l.Errorf(ctx, "cart.delivery.http.Get: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, cart)
}
