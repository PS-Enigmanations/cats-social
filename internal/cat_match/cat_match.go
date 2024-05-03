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
    CatMatchId               int      `json:"cat_match_id"`
    Message                  string   `json:"message"`
    Status                   string   `json:"status"`
    MatchCreatedAt           time.Time `json:"match_created_at"`
    UserName                 string   `json:"user_name"`
    UserEmail                string   `json:"user_email"`
    UserCreatedAt            time.Time `json:"user_created_at"`
    UserCatId                int      `json:"user_cat_id"`
    UserCatUserId            int      `json:"user_cat_user_id"`
    UserCatName              string   `json:"user_cat_name"`
    UserCatRace              string   `json:"user_cat_race"`
    UserCatSex               string   `json:"user_cat_sex"`
    UserCatAgeInMonth        int      `json:"user_cat_age_in_month"`
    UserCatDescription       string   `json:"user_cat_description"`
    UserCatIsAlreadyMatched  bool     `json:"user_cat_has_matched"`
    UserCatCreatedAt         time.Time `json:"user_cat_created_at"`
    UserCatUpdatedAt         time.Time `json:"user_cat_updated_at"`
    MatchCatId               int      `json:"match_cat_id"`
    MatchCatUserId           int      `json:"match_cat_user_id"`
    MatchCatName             string   `json:"match_cat_name"`
    MatchCatRace             string   `json:"match_cat_race"`
    MatchCatSex              string   `json:"match_cat_sex"`
    MatchCatAgeInMonth       int      `json:"match_cat_age_in_month"`
    MatchCatDescription      string   `json:"match_cat_description"`
    MatchCatIsAlreadyMatched bool     `json:"match_cat_has_matched"`
    MatchCatCreatedAt        time.Time `json:"match_cat_created_at"`
    MatchCatUpdatedAt        time.Time `json:"match_cat_updated_at"`
    UserCatImageUrls         []string `json:"user_cat_image_urls"`
    MatchCatImageUrls        []string `json:"match_cat_image_urls"`
}

type Status string

const (
	Pending  Status = "pending"
	Rejected Status = "rejected"
)

type UserDetail struct {
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"createdAt"`
}

type CatDetail struct {
    Id          int       `json:"id"`
    Name        string    `json:"name"`
    Race        string    `json:"race"`
    Sex         string    `json:"sex"`
    Description string    `json:"description"`
    AgeInMonth  int       `json:"ageInMonth"`
    ImageUrls   []string  `json:"imageUrls"`
    HasMatched  bool      `json:"hasMatched"`
    CreatedAt   time.Time `json:"createdAt"`
}

type CatMatchResponse struct {
    Id          int         `json:"id"`
    IssuedBy    UserDetail  `json:"issuedBy"`
    MatchCatDetail  CatDetail `json:"matchCatDetail"`
    UserCatDetail  CatDetail `json:"userCatDetail"`
    Message     string      `json:"message"`
    CreatedAt   time.Time   `json:"createdAt"`
}

func CatToCatResponse(c CatMatchValue) *CatMatchResponse {
	return &CatMatchResponse{
        Id:        c.CatMatchId,
        IssuedBy: UserDetail{
            Name:      c.UserName,
            Email:     c.UserEmail,
            CreatedAt: c.UserCreatedAt,
        },
        MatchCatDetail: CatDetail{
            Id:          c.MatchCatId,
            Name:        c.MatchCatName,
            Race:        c.MatchCatRace,
            Sex:         c.MatchCatSex,
            Description: c.MatchCatDescription,
            AgeInMonth:  c.MatchCatAgeInMonth,
            ImageUrls:   c.MatchCatImageUrls,
            HasMatched:  c.MatchCatIsAlreadyMatched,
            CreatedAt:   c.MatchCatCreatedAt,
        },
        UserCatDetail: CatDetail{
            Id:          c.UserCatId,
            Name:        c.UserCatName,
            Race:        c.UserCatRace,
            Sex:         c.UserCatSex,
            Description: c.UserCatDescription,
            AgeInMonth:  c.UserCatAgeInMonth,
            ImageUrls:   c.UserCatImageUrls,
            HasMatched:  c.UserCatIsAlreadyMatched,
            CreatedAt:   c.UserCatCreatedAt,
        },
        Message:   c.Message,
        CreatedAt: c.MatchCreatedAt,
    }
}