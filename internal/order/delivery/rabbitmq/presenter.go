package rabbitmq

type OrderMessage struct {
	CheckoutID string `json:"checkout_id"`
	OrderID    string `json:"order_id"`
	UserID     string `json:"user_id"`
}
