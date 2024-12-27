package resources

import "time"

type OrderEmailData struct {
	OrderID       string
	OrderDate     time.Time
	PaymentMethod string
	CustomerName  string
	CustomerEmail string
	Products      []OrderProductEmail
	TotalPrice    float64
	ShopTotal     float64
}

type OrderProductEmail struct {
	ProductID   string
	ProductName string
	Quantity    int
	Price       float32
	ShopName    string
	SubTotal    float64
	ShopID      string
}
