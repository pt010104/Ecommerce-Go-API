package http

import (
	"time"

	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/shop"
	"github.com/pt010104/api-golang/pkg/mongo"
	"github.com/pt010104/api-golang/pkg/paginator"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	City     string `json:"city" binding:"required"`
	Street   string `json:"street" binding:"required"`
	District string `json:"district"`
}

func (r registerRequest) validate() error {
	return nil
}

func (r registerRequest) toInput() shop.CreateShop {
	return shop.CreateShop{
		Name:     r.Name,
		Phone:    r.Phone,
		City:     r.City,
		Street:   r.Street,
		District: r.District,
	}
}

type address struct {
	City     string `json:"city"`
	Street   string `json:"street"`
	District string `json:"district,omitempty"`
}

type registerResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address address `json:"address"`
	Phone   string  `json:"phone"`
}

func (h handler) newCreateResponse(s models.Shop) registerResponse {
	return registerResponse{
		ID:    s.ID.Hex(),
		Name:  s.Name,
		Phone: s.Phone,
		Address: address{
			City:     s.City,
			Street:   s.Street,
			District: s.District,
		},
	}
}

type getShopRequest struct {
	IDs        []string `form:"ids"`
	Search     string   `form:"search"`
	IsVerified *bool    `form:"is_verified"`
}

func (r getShopRequest) validate() error {
	for _, id := range r.IDs {
		if !mongo.IsObjectID(id) {
			return errWrongBody
		}
	}

	return nil
}

func (r getShopRequest) toInput() shop.GetShopsFilter {
	return shop.GetShopsFilter{
		IDs:        r.IDs,
		Search:     r.Search,
		IsVerified: r.IsVerified,
	}
}

type listMetaResponse struct {
	paginator.PaginatorResponse
}

type getShopRespItem struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Address    address   `json:"address"`
	Followers  []string  `json:"followers,omitempty"`
	AvgRate    float64   `json:"avg_rate"`
	IsVerified *bool     `json:"is_verified,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type getShopResp struct {
	Meta  listMetaResponse  `json:"meta"`
	Items []getShopRespItem `json:"items"`
}

func (h handler) newGetShopsResp(ucOutput shop.GetShopOutput) getShopResp {
	var items []getShopRespItem
	for _, s := range ucOutput.Shops {
		shopItem := getShopRespItem{
			ID:    s.ID.Hex(),
			Name:  s.Name,
			Phone: s.Phone,
			Address: address{
				City:     s.City,
				Street:   s.Street,
				District: s.District,
			},
			AvgRate:    s.AvgRate,
			Followers:  mongo.HexFromObjectIDsOrNil(s.Followers),
			UserID:     s.UserID.Hex(),
			IsVerified: &s.IsVerified,
			CreatedAt:  s.CreatedAt,
		}

		items = append(items, shopItem)
	}

	return getShopResp{
		Meta: listMetaResponse{
			PaginatorResponse: ucOutput.Pag.ToResponse(),
		},
		Items: items,
	}
}

type getDetailResp struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Address    address   `json:"address"`
	Followers  []string  `json:"followers,omitempty"`
	AvgRate    float64   `json:"avg_rate"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

func (h handler) newDetailResponse(s models.Shop) getDetailResp {
	return getDetailResp{
		ID:    s.ID.Hex(),
		Name:  s.Name,
		Phone: s.Phone,
		Address: address{
			City:     s.City,
			Street:   s.Street,
			District: s.District,
		},
		AvgRate:    s.AvgRate,
		Followers:  mongo.HexFromObjectIDsOrNil(s.Followers),
		IsVerified: s.IsVerified,
		CreatedAt:  s.CreatedAt,
	}
}

type updateResp struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Address address `json:"address"`
	AvgRate float64 `json:"avg_rate"`
}
type updateShopRequest struct {
	IDs      []string `json:"ids" binding:"required"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	City     string   `json:"city"`
	Street   string   `json:"street"`
	District string   `json:"district"`
}

func (r updateShopRequest) toInput() shop.UpdateInput {
	return shop.UpdateInput{
		ShopIDs:  r.IDs,
		Name:     r.Name,
		Phone:    r.Phone,
		City:     r.City,
		Street:   r.Street,
		District: r.District,
	}
}
func (h handler) newUpdateShopResp(shops []models.Shop) []updateResp {
	var res []updateResp
	for _, s := range shops {
		res = append(res, updateResp{
			ID:    s.ID.Hex(),
			Name:  s.Name,
			Phone: s.Phone,
			Address: address{
				City:     s.City,
				Street:   s.Street,
				District: s.District,
			},
			AvgRate: s.AvgRate,
		})
	}

	return res
}
