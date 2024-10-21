package http

import (
	"github.com/pt010104/api-golang/internal/inventory"
	"github.com/pt010104/api-golang/pkg/mongo"
)

type createReq struct {
	ProductID       string `json:"product_id" binding:"required"`
	StockLevel      int    `json:"stock_level" binding:"required"`
	ReorderLevel    *int   `json:"reorder_level"`
	ReorderQuantity *int   `json:"reorder_quantity"`
}

func (r createReq) validate() error {
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

func (r createReq) toInput() inventory.CreateInput {
	input := inventory.CreateInput{
		ProductID:  r.ProductID,
		StockLevel: r.StockLevel,
	}

	if r.ReorderLevel != nil {
		input.ReorderLevel = r.ReorderLevel

		input.ReorderQuantity = r.ReorderQuantity
	}

	return input
}

type createResp struct {
	ID              string `json:"id"`
	ProductID       string `json:"product_id"`
	StockLevel      int    `json:"stock_level"`
	ReorderLevel    *int   `json:"reorder_level"`
	ReorderQuantity *int   `json:"reorder_quantity"`
}

func (h handler) newCreateResp(output inventory.CreateOutput) createResp {
	resp := createResp{
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
