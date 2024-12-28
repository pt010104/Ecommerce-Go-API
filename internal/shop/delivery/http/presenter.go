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

type Avatar_obj struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
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
	if len(r.IDs) > 0 {
		for _, id := range r.IDs {
			if !mongo.IsObjectID(id) {
				return errWrongBody
			}
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
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Name       string      `json:"name"`
	Phone      string      `json:"phone"`
	Address    address     `json:"address"`
	Followers  []string    `json:"followers,omitempty"`
	AvgRate    float64     `json:"avg_rate"`
	IsVerified *bool       `json:"is_verified,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	Avatar_obj *Avatar_obj `json:"avatar_obj,omitempty"`
}

type getShopResp struct {
	Meta  listMetaResponse  `json:"meta"`
	Items []getShopRespItem `json:"items"`
}

func (h handler) newGetShopsResp(ucOutput shop.GetShopOutput) getShopResp {
	var items []getShopRespItem
	for _, s := range ucOutput.Shops {
		shopItem := getShopRespItem{
			ID:    s.Shop.ID.Hex(),
			Name:  s.Shop.Name,
			Phone: s.Shop.Phone,
			Address: address{
				City:     s.Shop.City,
				Street:   s.Shop.Street,
				District: s.Shop.District,
			},
			AvgRate:    s.Shop.AvgRate,
			Followers:  mongo.HexFromObjectIDsOrNil(s.Shop.Followers),
			UserID:     s.Shop.UserID.Hex(),
			IsVerified: &s.Shop.IsVerified,
			CreatedAt:  s.Shop.CreatedAt,
		}

		if s.Avatar.URL != "" {
			shopItem.Avatar_obj = &Avatar_obj{
				MediaID: s.Avatar.MediaID,
				URL:     s.Avatar.URL,
			}
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
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	Address    address    `json:"address"`
	Followers  []string   `json:"followers,omitempty"`
	AvgRate    float64    `json:"avg_rate"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at"`
	Avatar_obj Avatar_obj `json:"avatar_obj"`
}

func (h handler) newDetailResponse(s shop.DetailShopOutput) getDetailResp {
	return getDetailResp{
		ID:    s.S.ID.Hex(),
		Name:  s.S.Name,
		Phone: s.S.Phone,
		Address: address{
			City:     s.S.City,
			Street:   s.S.Street,
			District: s.S.District,
		},
		AvgRate:    s.S.AvgRate,
		Followers:  mongo.HexFromObjectIDsOrNil(s.S.Followers),
		IsVerified: s.S.IsVerified,
		CreatedAt:  s.S.CreatedAt,
		Avatar_obj: Avatar_obj{
			MediaID: s.MediaID,
			URL:     s.URL,
		},
	}
}

type updateResp struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Address address `json:"address"`
	AvgRate float64 `json:"avg_rate"`
}
type GetShopIDByUserIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (r GetShopIDByUserIDRequest) validate() error {

	if !mongo.IsObjectID(r.ID) {
		return errWrongBody
	}

	return nil
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

type MostSolProductItem struct {
	ProductID   string `json:"id"`
	ProductName string `json:"name"`
	Sold        int    `json:"sold"`
}

type MostViewedProductItem struct {
	ProductID   string `json:"id"`
	ProductName string `json:"name"`
	Viewed      int    `json:"viewed"`
}

type MostViewTrendItem struct {
	ProductID   string `json:"id"`
	ProductName string `json:"name"`
	ViewTrend   int    `json:"view_trend"`
}

type reportResponse struct {
	MostViewedProducts []MostViewedProductItem `json:"most_viewed_products"`
	MostSoldProducts   []MostSolProductItem    `json:"most_sold_products"`
	MostViewTrend      []MostViewTrendItem     `json:"most_view_trend"`
}

func (h handler) newReportResponse(report shop.ReportOutput) reportResponse {
	mostViewedProducts := make([]MostViewedProductItem, len(report.MostViewedProducts))
	for i, v := range report.MostViewedProducts {
		mostViewedProducts[i] = MostViewedProductItem{
			ProductID:   v.ID.Hex(),
			ProductName: v.Name,
			Viewed:      v.View,
		}
	}
	mostSoldProducts := make([]MostSolProductItem, len(report.MostSoldProducts))
	for i, v := range report.MostSoldProducts {
		mostSoldProducts[i] = MostSolProductItem{
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			Sold:        v.Sold,
		}
	}

	mostViewTrend := make([]MostViewTrendItem, len(report.MostViewTrend))
	for i, v := range report.MostViewTrend {
		mostViewTrend[i] = MostViewTrendItem{
			ProductID:   v.ID.Hex(),
			ProductName: v.Name,
			ViewTrend:   v.ViewTrend,
		}
	}

	return reportResponse{
		MostViewedProducts: mostViewedProducts,
		MostSoldProducts:   mostSoldProducts,
		MostViewTrend:      mostViewTrend,
	}
}
