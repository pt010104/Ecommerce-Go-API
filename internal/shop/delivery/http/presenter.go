package http

import "github.com/pt010104/api-golang/internal/shop"

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	City     string `json:"city" binding:"required"`
	Street   string `json:"street" binding:"required"`
	District string `json:"district" binding:"required"`
}

func (r registerRequest) validate() error {
	return nil
}

func (r registerRequest) toInput() shop.RegisterInput {
	return shop.RegisterInput{
		Name:     r.Name,
		City:     r.City,
		Street:   r.Street,
		District: r.District,
	}
}
