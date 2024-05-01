package catmatch

import (
	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/user"
	"time"
)

type CatMatch struct {
	Id       int    `json:"id"`
	IssuedBy int64  `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Cat  cat.Cat
	User user.User

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
