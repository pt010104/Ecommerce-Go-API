package order

import (
	"time"

	"github.com/pt010104/api-golang/internal/models"
)

type CreateCheckoutOutput struct {
	CheckoutID       string             `json:"checkout_id"`
	ExpiredAt        time.Time          `json:"expired_at"`
	TotalPriceByShop map[string]float64 `json:"total_price_by_shop"`
	TotalPrice       float64            `json:"total_price"`
	Products         []models.Product   `json:"products"`
	QuantityMap      map[string]int     `json:"quantity_map"`
	Shops            []models.Shop      `json:"shops"`
	ImageURLs        []string           `json:"image_urls"`
}
