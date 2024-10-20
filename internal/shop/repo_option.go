package shop

import (
	"github.com/pt010104/api-golang/pkg/paginator"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateOption struct {
	Name     string
	Alias    string
	City     string
	Street   string
	District string
	Phone    string
}

type GetOption struct {
	GetShopsFilter
	PagQuery paginator.PaginatorQuery
}

type DetailOption struct {
	ID string
}
type UpdateOption struct {
	ID         string
	UpdateData bson.M
}
