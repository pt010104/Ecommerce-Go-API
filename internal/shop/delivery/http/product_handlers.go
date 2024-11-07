package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h handler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processCreateProductRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Create: %v", err)
		response.Error(c, err)
		return
	}

	product, _, err := h.uc.CreateProduct(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Create: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, product)
}
func (h handler) DetailProduct(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processDetailProductRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.detail: %v", err)
		response.Error(c, err)
		return
	}
	productID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		h.l.Errorf(ctx, "invalid product ID format: %v", err)
		response.Error(c, err)
		return
	}
	fmt.Println("ID:", req.ID)
	product, inven, err := h.uc.DetailProduct(ctx, sc, productID)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.detauk: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newDetailProductResponse(product, inven))
}
