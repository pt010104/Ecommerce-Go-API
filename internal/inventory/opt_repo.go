package inventory

type CreateOption struct {
	ProductID       string
	StockLevel      int
	ReorderLevel    *int
	ReorderQuantity *int
}
