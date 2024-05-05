package session

import "time"

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
