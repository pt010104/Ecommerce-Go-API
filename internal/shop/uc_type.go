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
	Name      *string               `json:"name,omitempty"`
	Alias     *string               `json:"alias,omitempty"`
	City      *string               `json:"city,omitempty"`
	Street    *string               `json:"street,omitempty"`
	District  *string               `json:"district,omitempty"`
	Phone     *string               `json:"phone,omitempty"`
	Followers *[]primitive.ObjectID `json:"followers,omitempty"`
	AvgRate   *float64              `json:"avg_rate,omitempty"`
}
type GetOutput struct {
	Shops []models.Shop
	Pag   paginator.Paginator
}
