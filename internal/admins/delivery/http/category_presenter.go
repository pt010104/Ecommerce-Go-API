package http

import (
	"github.com/pt010104/api-golang/internal/admins"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/mongo"
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

func (h handler) newCategoryResponse(u models.Category) CategoryResponse {
	return CategoryResponse{
		ID:          u.ID.Hex(),
		Name:        u.Name,
		Description: u.Description,
	}
}

type listCatagoryReq struct {
	Name string `form:"name"`

	IDs []string `form:"ids"`
}

func (r listCatagoryReq) toInput() admins.GetCategoriesFilter {

	return admins.GetCategoriesFilter{
		IDs:  r.IDs,
		Name: r.Name,
	}
}
func (r listCatagoryReq) validate() error {
	for _, id := range r.IDs {
		if !mongo.IsObjectID(id) {
			return errWrongBody
		}
	}
	return nil
}

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type listCategoryResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

func (h handler) newListResponse(categories []models.Category) listCategoryResponse {
	var resList []CategoryResponse
	for _, category := range categories {
		res := CategoryResponse{
			ID:          category.ID.Hex(),
			Name:        category.Name,
			Description: category.Description,
		}
		resList = append(resList, res)
	}
	return listCategoryResponse{
		Categories: resList,
	}
}
