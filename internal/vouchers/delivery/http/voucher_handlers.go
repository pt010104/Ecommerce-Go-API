package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary Create a voucher
// @Schemes http https
// @Description Create a new voucher
// @Tags Voucher
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param session-id header string true "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param request body createVoucherReq true "Voucher creation request"
// @Success 200 {object} response.Resp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/vouchers [post]
func (h handler) CreateVoucher(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateVoucherRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.processCreateRequest: %v", err)
		response.Error(c, err)
		return
	}
	u, err := h.uc.CreateVoucher(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, u)
}

// @Summary Get voucher details
// @Schemes http https
// @Description Get details of a specific voucher , pass id or code as param one of them must be presented if search by id change the route to by-id/id
// @Tags Voucher
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param session-id header string true "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param data body DetailVoucherReq true "Detail Voucher Request"
// @Success 200 {object} response.Resp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/vouchers/by-code{code} [get]
func (h handler) DetailVoucher(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processDetailVoucherRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.processDetailRequest: %v", err)
		response.Error(c, err)
		return
	}
	u, err := h.uc.Detail(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newDetailResponse(u))
}

// @Summary List vouchers
// @Schemes http https
// @Description Get a list of vouchers with optional filtering
// @Tags Voucher
// @Accept json
// @Produce json
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
// @Param Authorization header string true "Bearer JWT token" default(Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjAxMTk2NjgsImlhdCI6MTcyODU4MzY2OCwic3ViIjoiNjcwNzgyNWQ0NTgwNGNhYWY4MzE2OTU3Iiwic2Vzc2lvbl9pZCI6InpnSFJMd1NmTnNQVnk2d2g3M0ZLVmpqZXV6T1ZnWGZSMjdRYVd1eGtsdzQ9IiwidHlwZSI6IiIsInJlZnJlc2giOmZhbHNlfQ.Pti0gJ5fO4WjGTsxShGv90pr0E_0jMJdWFEUJYKG4VU)
// @Param x-client-id header string true "User ID" default(6707825d45804caaf8316957)
// @Param session-id header string true "Session ID" default(zgHRLwSfNsPVy6wh73FKVjjeuzOVgXfR27QaWuxklw4=)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param shop_id query string false "Filter by shop ID"
// @Success 200 {object} response.Resp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
// @Router /api/v1/vouchers [get]
func (h handler) ListVoucher(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processListVoucherRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.processListRequest: %v", err)
		response.Error(c, err)
		return
	}
	u, err := h.uc.List(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newListResponse(u))
}
