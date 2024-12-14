package http

import (
	"fmt"

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

// @Summary		get shop detail by their id
// @Schemes		http https
// @Description	Get shop detail by id
// @Tags			Products
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			id							path		string		true	"User ID"
//
// @Success		200							{object}	detailProductResp	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops/products/{id} [GET]
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

// @Summary		delete product by their id
// @Schemes		http https
// @Description	delete product by id
// @Tags			Products
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			request body deleteProductRequest true "Request Body"
//
// @Success		200							{object}	response.Resp	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops/products/delete [DELETE]
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
// @param       request query getProductRequest false "Request Body"
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200 {object} getProductResp "Paginated Products"
// @Router       /api/v1/shops/products/get-product [get]
func (h handler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()
	fmt.Print("get product", c.Params, c.Query("ids"))
	sc, req, err := h.processGetProductRequest(c)
	if err != nil {
		fmt.Print("err", req.IDs)
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
		h.l.Errorf(ctx, "shop.delivery.http.getProduct: %v", err2)
		err := h.mapErrors(err)
		response.Error(c, err)
		return
	}
	response.OK(c, h.getProductResp(o))

}

// @Summary		Update product
// @Schemes		http https
// @Description	Update product BY ID , only ID is required , other fields are optional
// @Tags			Products
// @Accept			json
// @Produce		json
//
// @Param			Access-Control-Allow-Origin	header		string		false	"Access-Control-Allow-Origin"	default("*")
// @Param			Authorization				header		string		true	"Bearer JWT token"				default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param			x-client-id					header		string		true	"User ID"						default(6707825d45804caaf8316957)
// @Param			session-id					header		string		true	"Session ID"					default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param			request body UpdateProductReq true "Request Body"
//
// @Success		200							{object}	updateProductResp	"Success"
// @Failure		400							{object}	response.Resp	"Bad Request"
// @Failure		500							{object}	response.Resp	"Internal Server Error"
//
// @Router			/api/v1/shops/products/update [Post]
func (h handler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	sc, req, err := h.processUpdateProductRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "shop.delivery.http.UpdateProduct: %v", err)
		response.Error(c, err)
		return
	}

	p, err2 := h.uc.UpdateProduct(ctx, sc, req.toInput())
	if err2 != nil {
		h.l.Errorf(ctx, "shop.delivery.http.delete: %v", err2)
		err := h.mapErrors(err2)
		response.Error(c, err)
		return
	}
	response.OK(c, h.newUpdateProductResponse(p))

}
