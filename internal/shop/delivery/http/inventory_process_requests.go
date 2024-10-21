package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
)

func (h handler) processCreateInventoryRequest(c *gin.Context) (createInventoryReq, models.Scope, error) {
	ctx := c.Request.Context()

	var req createInventoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "inventory.delivery.http.handler.processCreateRequest: invalid request")
		return createInventoryReq{}, models.Scope{}, errWrongBody
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "inventory.delivery.http.handler.processCreateRequest: invalid request")
		return createInventoryReq{}, models.Scope{}, errWrongBody
	}

	return req, models.Scope{}, nil
}
