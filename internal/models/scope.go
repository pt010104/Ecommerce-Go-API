package models

type Scope struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	ShopID    string `json:"shop_id"`
}
