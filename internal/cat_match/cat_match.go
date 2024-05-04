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

type CatMatchValue struct {
	CatMatchId          int       `json:"cat_match_id"`
	Message             string    `json:"message"`
	Status              string    `json:"status"`
	MatchCreatedAt      time.Time `json:"match_created_at"`
	UserName            string    `json:"user_name"`
	UserEmail           string    `json:"user_email"`
	UserCreatedAt       time.Time `json:"user_created_at"`
	UserCatId           int       `json:"user_cat_id"`
	UserCatUserId       int       `json:"user_cat_user_id"`
	UserCatName         string    `json:"user_cat_name"`
	UserCatRace         string    `json:"user_cat_race"`
	UserCatSex          string    `json:"user_cat_sex"`
	UserCatAgeInMonth   int       `json:"user_cat_age_in_month"`
	UserCatDescription  string    `json:"user_cat_description"`
	UserCatHasMatched   bool      `json:"user_cat_has_matched"`
	UserCatCreatedAt    time.Time `json:"user_cat_created_at"`
	MatchCatId          int       `json:"match_cat_id"`
	MatchCatUserId      int       `json:"match_cat_user_id"`
	MatchCatName        string    `json:"match_cat_name"`
	MatchCatRace        string    `json:"match_cat_race"`
	MatchCatSex         string    `json:"match_cat_sex"`
	MatchCatAgeInMonth  int       `json:"match_cat_age_in_month"`
	MatchCatDescription string    `json:"match_cat_description"`
	MatchCatHasMatched  bool      `json:"match_cat_has_matched"`
	MatchCatCreatedAt   time.Time `json:"match_cat_created_at"`
	UserCatImageUrls    []string  `json:"user_cat_image_urls"`
	MatchCatImageUrls   []string  `json:"match_cat_image_urls"`
}

type Status string

const (
	Pending  Status = "pending"
	Rejected Status = "rejected"
)
