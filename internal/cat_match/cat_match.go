package catmatch

import (
	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/internal/user"
	"time"
)

type CatMatch struct {
	Id         int    `json:"id"`
	IssuedBy   int64  `json:"issued_by" validate:"required"`
	MatchCatId int64  `json:"match_cat_id" validate:"required"`
	UserCatId  int64  `json:"user_cat_id" validate:"required"`
	Message    string `json:"message"`
	Status     Status `json:"status" validate:"required,oneof=pending rejected"`

	Cat  cat.Cat
	User user.User

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type Status string

const (
	Pending  Status = "pending"
	Rejected Status = "rejected"
)
