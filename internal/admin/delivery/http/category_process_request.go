package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
)

func (h handler) processCreateCategoryRequest(c *gin.Context) (createCategoryReq, models.Scope, error) {
	ctx := c.Request.Context()
	var req createCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.processCreateRequest: invalid request")
		return createCategoryReq{}, models.Scope{}, errWrongInput
	}

	return req, models.Scope{}, nil

}
