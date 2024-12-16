package http

import (
	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/models"
)

type createCategoryReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (req createCategoryReq) toInput() admins.CreateCategoryInput {
	return admins.CreateCategoryInput{
		Name:        req.Name,
		Description: req.Description,
	}

}

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h handler) newCategoryResponse(u models.Category) CategoryResponse {
	return CategoryResponse{
		ID:          u.ID.Hex(),
		Name:        u.Name,
		Description: u.Description,
	}
}
