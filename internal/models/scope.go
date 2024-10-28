package models

type Scope struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Role      int
}

var (
	RoleUser  = 0
	RoleAdmin = 1
)
