package response

import "enigmanations/cats-social/internal/cat"

type CatMatchCreateResponse struct {
	Message string `json:"message"`
}

const CatMatchCreateSuccMessage = "Successfully send match request"

func CatToCatMatchCreateResponse(c cat.Cat) *CatMatchCreateResponse {
	return &CatMatchCreateResponse{
		Message: CatMatchCreateSuccMessage,
	}
}
