package response

import (
	"enigmanations/cats-social/internal/cat"
	catmatch "enigmanations/cats-social/internal/cat_match"
	"time"
)

// Create response
type CatMatchCreateResponse struct {
	Message string `json:"message"`
}

const CatMatchCreateSuccMessage = "Successfully send match request"

func CatToCatMatchCreateResponse(c cat.Cat) *CatMatchCreateResponse {
	return &CatMatchCreateResponse{
		Message: CatMatchCreateSuccMessage,
	}
}

// Get all response
type issuedBy struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type catDetail struct {
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

type CatGetAllShow struct {
	Id             int       `json:"id"`
	IssuedBy       issuedBy  `json:"issuedBy"`
	MatchCatDetail catDetail `json:"matchCatDetail"`
	UserCatDetail  catDetail `json:"userCatDetail"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"createdAt"`
}

type CatGetAllShows []CatGetAllShow

type CatGetAllResponse struct {
	Message string         `json:"message"`
	Data    CatGetAllShows `json:"data"`
}

const CatMatchGetAllSuccMessage = "Successfully get match request"

func CatToCatGetAllResponse(data CatGetAllShows) *CatGetAllResponse {
	return &CatGetAllResponse{
		Message: CatMatchGetAllSuccMessage,
		Data:    data,
	}
}

func ToCatMatchShows(cms []*catmatch.CatMatchValue) CatGetAllShows {
	list := make(CatGetAllShows, len(cms))

	for i, item := range cms {
		list[i] = CatGetAllShow{
			Id: item.CatMatchId,
			IssuedBy: issuedBy{
				Name:      item.UserName,
				Email:     item.UserEmail,
				CreatedAt: item.UserCreatedAt,
			},
			MatchCatDetail: catDetail{
				Id:          item.MatchCatId,
				Name:        item.MatchCatName,
				Race:        item.MatchCatRace,
				Sex:         item.MatchCatSex,
				Description: item.MatchCatDescription,
				AgeInMonth:  item.MatchCatAgeInMonth,
				ImageUrls:   item.MatchCatImageUrls,
				HasMatched:  item.MatchCatHasMatched,
				CreatedAt:   item.MatchCatCreatedAt,
			},
			UserCatDetail: catDetail{
				Id:          item.UserCatId,
				Name:        item.UserCatName,
				Race:        item.UserCatRace,
				Sex:         item.UserCatSex,
				Description: item.UserCatDescription,
				AgeInMonth:  item.UserCatAgeInMonth,
				ImageUrls:   item.UserCatImageUrls,
				HasMatched:  item.MatchCatHasMatched,
				CreatedAt:   item.UserCatCreatedAt,
			},
			Message:   item.Message,
			CreatedAt: item.MatchCreatedAt,
		}
	}

	return list
}
