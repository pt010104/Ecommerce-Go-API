package http

import "github.com/pt010104/api-golang/internal/admin"

type createCategoryReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (req createCategoryReq) toInput() admin.CreateCategoryInput {
	return admin.CreateCategoryInput{
		Name:        req.Name,
		Description: req.Description,
	}

}
