package request

type CatImageCreateRequest struct {
	ImageUrls []string `validate:"required,min=1,dive,required,url"`
}
