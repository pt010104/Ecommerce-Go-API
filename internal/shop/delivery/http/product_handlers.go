package http

import (
	"github.com/gin-gonic/gin"

	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/paginator"
	"github.com/pt010104/api-golang/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary      Create a product
// @Description  Create a new product in the shop
// @Tags         Shop
// @Accept       json
// @Produce      json
// @Param        request body createProductReq true "Request Body"
// @Success      200 {object} response.Resp "Success Response"
// @Failure      400 {object} response.Resp "Bad Request"
// @Failure      500 {object} response.Resp "Internal Server Error"
// @Router       /api/v1/shops/create-product [post]
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

// @Summary      Get product details
// @Description  Retrieve detailed information about a product by ID
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        request body detailProductReq true "Request Body"
// @Success      200 {object} detailProductResp "Product Details"
// @Failure      400 {object} response.Resp "Bad Request"
// @Router       /api/v1/shops/ [get]
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

	o, err2 := h.uc.ListProduct(ctx, sc, req.toInput())
	if err2 != nil {
		h.l.Errorf(ctx, "shop.delivery.http.listProduct: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}
	response.OK(c, h.listProductResp(o))

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
		err := h.mapErrors(err2)
		response.Error(c, err)
		return
	}
	response.OK(c, " delete success")

}

// @Summary      Get products with pagination
// @Description  Retrieve a paginated list of products with optional filters , at least 1 parameter request body must be provided
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        request body getProductRequest false "Request Body"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200 {object} getProductResp "Paginated Products"
// @Router       /api/v1/shops/products/get-product [get]
func (h handler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processGetProductRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.getproduct: %v", err)
		response.Error(c, err)
		return
	}
	var pagQuery paginator.PaginatorQuery
	if err := c.ShouldBindQuery(&pagQuery); err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.Get.ShouldBindQuery: %v", err)
		response.Error(c, errWrongPaginationQuery)
		return
	}

	o, err2 := h.uc.GetProduct(ctx, sc, shop.GetProductOption{
		GetProductFilter: req.toInput(),
		PagQuery:         pagQuery,
	})
	if err2 != nil {
		h.l.Errorf(ctx, "shop.delivery.http.getProduct: %v", err)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}
	response.OK(c, h.getProductResp(o))

}
