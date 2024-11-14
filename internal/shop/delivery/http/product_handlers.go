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
	product, err2 := h.uc.DetailProduct(ctx, sc, productID)
	if err2 != nil {
		h.l.Errorf(ctx, "shop.delivery.http.detauk: %v", err)
		err2 := h.mapErrors(err2)
		response.Error(c, err2)
		return
	}

	response.OK(c, h.newDetailProductResponse(product))
}
func (h handler) ListProduct(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processListProductRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.listproduct: %v", err)
		response.Error(c, err)
		return
	}

	list, err2 := h.uc.ListProduct(ctx, sc, req.toInput())
	if err2 != nil {
		h.l.Errorf(ctx, "shop.delivery.http.detauk: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}
	response.OK(c, h.listProductResp(list.List))

}
func (h handler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processDeleteProductRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.DeleteProduct: %v", err)
		response.Error(c, err)
		return
	}

	err2 := h.uc.DeleteProduct(ctx, sc, req.IDs)
	if err2 != nil {
		h.l.Errorf(ctx, "shop.delivery.http.delete: %v", err2)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}
	response.OK(c, " delete success")

}
