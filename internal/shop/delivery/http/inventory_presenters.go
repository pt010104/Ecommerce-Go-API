package http

import (
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
)

type createInventoryReq struct {
	ProductID       string `json:"product_id" binding:"required"`
	StockLevel      int    `json:"stock_level" binding:"required"`
	ReorderLevel    *int   `json:"reorder_level"`
	ReorderQuantity *int   `json:"reorder_quantity"`
}

func (r createInventoryReq) validate() error {
	if r.ReorderLevel != nil {
		if *r.ReorderLevel < 0 {
			return errWrongBody
		}

		if r.ReorderQuantity == nil {
			return errWrongBody
		}

		if r.ReorderQuantity != nil {
			if *r.ReorderQuantity < 0 {
				return errWrongBody
			}
		}
	}

	if !mongo.IsObjectID(r.ProductID) {
		return errWrongBody
	}

	return nil
}

func (r createInventoryReq) toInput() shop.CreateInventoryInput {
	input := shop.CreateInventoryInput{
		ProductID:  r.ProductID,
		StockLevel: r.StockLevel,
	}

	if r.ReorderLevel != nil {
		input.ReorderLevel = r.ReorderLevel

		input.ReorderQuantity = r.ReorderQuantity
	}

	return input
}

type createInventoryResp struct {
	ID              string `json:"id"`
	ProductID       string `json:"product_id"`
	StockLevel      int    `json:"stock_level"`
	ReorderLevel    *int   `json:"reorder_level,omitempty"`
	ReorderQuantity *int   `json:"reorder_quantity,omitempty"`
}

func (h handler) newCreateInventoryResp(output shop.CreateInventoryOutput) createInventoryResp {
	resp := createInventoryResp{
		ID:         output.Inventory.ID.Hex(),
		ProductID:  output.Inventory.ProductID.Hex(),
		StockLevel: output.Inventory.StockLevel,
	}

	if output.Inventory.ReorderLevel != nil {
		resp.ReorderLevel = output.Inventory.ReorderLevel
		resp.ReorderQuantity = output.Inventory.ReorderQuantity
	}

	return resp
}
