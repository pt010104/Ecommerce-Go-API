package http

import (
	"github.com/pt010104/api-golang/internal/shop"
)

type createProductReq struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`

	StockLevel      uint  `json:"stock_level"`
	ReorderLevel    *uint `json:"reorder_level"`
	ReorderQuantity *uint `json:"reorder_quantity"`
}

func (r createProductReq) toInput() shop.CreateProductInput {
	return shop.CreateProductInput{
		Name:  r.Name,
		Price: r.Price,

		StockLevel:      r.StockLevel,
		ReorderLevel:    r.ReorderLevel,
		ReorderQuantity: r.ReorderQuantity,
	}
}
