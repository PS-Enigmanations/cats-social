package web

type CatCreateRequest struct {
	Name        string `validate:"required, min=1, max=30"`
	Race        string `validate:"required"`
	Sex         string `validate:"required"`
	AgeInMonth  int    `validate:"required"`
	Description string
	ImageUrls   string
}
