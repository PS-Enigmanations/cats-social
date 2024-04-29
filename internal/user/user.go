package user

import (
	"time"

	"github.com/gofrs/uuid"
)

type UserSession struct {
	AccessToken  string    `json:"access_token"`
	TokenExpires int       `json:"token_expires"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type User struct {
	ID       string    `json:"id"`
	UUID     uuid.UUID `json:"uid"`
	Email    *string   `json:"email"`
	Password string    `json:"password"`
	Phone    *string   `json:"phone"`

	Session []UserSession

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
