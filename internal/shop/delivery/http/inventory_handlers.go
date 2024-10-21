package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary Create an inventory product
// @Schemes http https
// @Description Create an inventory product
// @Tags Inventory
// @Accept json
// @Produce json
//
// @Param Access-Control-Allow-Origin header string false "Access-Control-Allow-Origin" default("*")
//
// @Success 200 {object} createInventoryResp
// @Failure 400 {object} response.Resp "Bad Request"
// @Failure 500 {object} response.Resp "Internal Server Error"
//
// @Router /api/v1/inventories/ [POST]
func (h handler) CreateInventory(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateInventoryRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "inventory.delivery.http.handler.Create.processCreateRequest: %v", err)
		response.Error(c, err)
		return
	}

	u, err := h.uc.CreateInventory(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "inventory.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}

	response.OK(c, h.newCreateInventoryResp(u))
}
