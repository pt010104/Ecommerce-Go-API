package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/response"
)

// @Summary Create Category
// @Tags Admin
// @Description Create a new category
// @Accept json
// @Produce json
// @Param category body createCategoryReq true "Category data"
// @Success 200 {object} CategoryResponse
// @Router /api/v1/admin/categories [post]
func (h handler) CreateCategory(c *gin.Context) {
	ctx := c.Request.Context()

	req, sc, err := h.processCreateCategoryRequest(c)
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.processCreateRequest: %v", err)
		response.Error(c, err)
		return
	}
	h.l.Debugf(ctx, "role:", sc.Role)
	u, err := h.uc.CreateCategory(ctx, sc, req.toInput())
	if err != nil {
		h.l.Errorf(ctx, "category.delivery.http.handler.Create.Create: %v", err)
		err = h.mapErrors(err)
		response.Error(c, err)
		return
	}
	h.l.Debugf(ctx, "role:", sc.Role)
	response.OK(c, h.newCategoryResponse(u))
}
