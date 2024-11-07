package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
)

type createProductReq struct {
	Name            string  `json:"name" binding:"required"`
	Price           float32 `json:"price" binding:"required"`
	StockLevel      uint    `json:"stock_level" binding:"required"`
	ReorderLevel    *uint   `json:"reorder_level" binding:"required"`
	ReorderQuantity *uint   `json:"reorder_quantity" binding:"required"`
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

type detailProductReq struct {
	ID string `json:"id" binding:"required"`
}
type detailProductResp struct {
	ID          string  `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	InventoryID string  `json:"inventory_id" binding:"required"`
	Price       float32 `json:"price" binding:"required"`
}

func (h handler) newDetailProductResponse(p models.Product, i models.Inventory) detailProductResp {
	return detailProductResp{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		InventoryID: i.ID.Hex(),
		Price:       p.Price,
	}

}
