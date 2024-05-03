package request

type CatMatchRequest struct {
	MatchCatId int64  `json:"matchCatId" validate:"required"`
	UserCatId  int64  `json:"userCatId" validate:"required"`
	Message    string `json:"message" validate:"required,min=5,max=120"`
}

type CatMatchApproveRejectRequest struct {
	MatchId int64  `json:"matchId" validate:"required"`
}
