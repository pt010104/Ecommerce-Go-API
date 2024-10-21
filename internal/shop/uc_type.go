package shop

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/paginator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateInput struct {
	Name     string
	City     string
	Street   string
	District string
	Phone    string
}

type GetShopsFilter struct {
	IDs    []string
	Search string
}
type DeleteInput struct {
	ID string
}

type GetInput struct {
	GetShopsFilter
	PagQuery paginator.PaginatorQuery
}
type UpdateInput struct {
	ID        string
	Name      *string
	Alias     *string
	City      *string
	Street    *string
	District  *string
	Phone     *string
	Followers *[]primitive.ObjectID
	AvgRate   *float64
}
type GetOutput struct {
	Shops []models.Shop
	Pag   paginator.Paginator
}
