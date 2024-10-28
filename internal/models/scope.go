package models

type Scope struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Role      int
}
