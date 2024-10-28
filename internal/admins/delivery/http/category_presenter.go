package http

import "github.com/pt010104/api-golang/internal/admins"

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
