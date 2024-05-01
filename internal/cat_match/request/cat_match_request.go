package request

type CatMatchRequest struct {
	MatchCatId int `json:"matchCatId" validate:"required"`
	UserCatId  int `json:"userCatId" validate:"required"`
}
