package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createProductReq struct {
	Name            string   `json:"name" binding:"required"`
	Price           float32  `json:"price" binding:"required"`
	StockLevel      uint     `json:"stock_level" binding:"required"`
	ReorderLevel    *uint    `json:"reorder_level" binding:"required"`
	ReorderQuantity *uint    `json:"reorder_quantity" binding:"required"`
	CategoryIDs     []string `json:"category_ids" binding:"required"`
}

func (r createProductReq) toInput() shop.CreateProductInput {
	return shop.CreateProductInput{
		Name:  r.Name,
		Price: r.Price,

		StockLevel:      r.StockLevel,
		ReorderLevel:    r.ReorderLevel,
		ReorderQuantity: r.ReorderQuantity,
		CategoryID:      r.CategoryIDs,
	}
}
func (r createProductReq) validate() error {

	if r.Name == "" {
		return errWrongBody
	}

	if r.Price <= 0 {
		return errWrongBody
	}

	if r.StockLevel == 0 {
		return errWrongBody
	}

	if r.ReorderLevel == nil || *r.ReorderLevel == 0 {
		return errWrongBody
	}

	if r.ReorderQuantity == nil || *r.ReorderQuantity == 0 {
		return errWrongBody
	}

	for _, id := range r.CategoryIDs {
		if _, err := primitive.ObjectIDFromHex(id); err != nil {
			return errWrongBody
		}
	}

	return nil
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

type listProductRequest struct {
	IDs    []string `json:"ids"`
	Search string   `json:"search"`
	ShopID []string `json:"shop_id"`
}

func (r listProductRequest) validate() error {
	if len(r.IDs) > 0 {
		for _, id := range r.IDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
		}
	}

	return nil
}

func (r listProductRequest) toInput() shop.GetProductFilter {
	return shop.GetProductFilter{
		IDs:    r.IDs,
		Search: r.Search,
		ShopID: r.ShopID,
	}
}

type listProductItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	ShopID      string  `json:"shop_id"`
	InventoryID string  `json:"inventory_id"`
	Price       float32 `json:"price"`
}

type listProductResp struct {
	Item []listProductItem
}

func (h handler) listProductResp(p []models.Product) listProductResp {
	var list []listProductItem
	for _, s := range p {
		item := listProductItem{
			ID:          s.ID.Hex(),
			Name:        s.Name,
			ShopID:      s.ShopID.Hex(),
			InventoryID: s.InventoryID.Hex(),
			Price:       s.Price,
		}
		list = append(list, item)
	}
	return listProductResp{
		Item: list,
	}

}
