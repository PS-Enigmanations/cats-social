package response

import (
	"enigmanations/cats-social/internal/cat"
	"enigmanations/cats-social/pkg/structure"
	"time"
)

// Get all response
type CatShow struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Race        string    `json:"race"`
	Sex         string    `json:"sex"`
	AgeInMonth  int       `json:"ageInMonth"`
	ImageUrls   []string  `json:"imageUrls"`
	HasMatched  bool      `json:"hasMatched"`
	Description string    `json:"description" validate:"required,min=1,max=200"`
	CreatedAt   time.Time `json:"created_at"`
}

type CatShows []CatShow

const CatGetAllSuccMessage = "Successfully get cats"

type CatGetAllResponse struct {
	Message string   `json:"message"`
	Data    CatShows `json:"data"`
}

func CatToCatGetAllResponse(data CatShows) *CatGetAllResponse {
	return &CatGetAllResponse{
		Message: CatGetAllSuccMessage,
		Data:    data,
	}
}

func ToCatShows(c []*cat.Cat) CatShows {
	list := make(CatShows, len(c))
	for i, item := range c {
		showItem := new(CatShow)
		structure.Copy(item, showItem)
		list[i] = *showItem
	}

	return list
}

// Create response
type CatCreateShow struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type CatCreateResponse struct {
	Message string        `json:"message"`
	Data    CatCreateShow `json:"data"`
}

const CatCreateSuccMessage = "Successfully add cat"

func CatToCatCreateResponse(c cat.Cat) *CatCreateResponse {
	return &CatCreateResponse{
		Message: CatCreateSuccMessage,
		Data: CatCreateShow{
			Id:        c.Id,
			CreatedAt: c.CreatedAt,
		},
	}
}
