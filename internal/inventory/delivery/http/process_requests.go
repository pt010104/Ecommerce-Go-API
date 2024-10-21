package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
)

func (h handler) processCreateRequest(c *gin.Context) (createReq, models.Scope, error) {
	ctx := c.Request.Context()

	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "inventory.delivery.http.handler.processCreateRequest: invalid request")
		return createReq{}, models.Scope{}, err
	}

	if err := req.validate(); err != nil {
		h.l.Errorf(ctx, "inventory.delivery.http.handler.processCreateRequest: invalid request")
		return createReq{}, models.Scope{}, errWrongBody
	}

	return req, models.Scope{}, nil
}
