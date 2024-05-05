package catimage

type CatImage struct {
	Id    int    `json:"id"`
	CatId int    `json:"catId"`
	Url   string `json:"url" validate:"required,url"`
}
