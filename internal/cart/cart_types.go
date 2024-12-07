package cart

import ()

type CreateCartInput struct {
	UserID string
	ShopID string
	Item   []CreateCartItemInput
}
type CreateCartItemInput struct {
	ProductID string
	Quantity  int
}
type CreateCartOutput struct {
	ID     string
	UserID string
	ShopID string
	Item   []CreateCartItemOutput
}
type CreateCartItemOutput struct {
	ProductID string
	Quantity  int
}
