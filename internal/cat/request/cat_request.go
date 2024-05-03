package request

type CatGetAllRequestParams struct {
	Id         string `json:"id"`
	Limit      string `json:"limit" default:"5"`
	Offset     string `json:"offset" default:"0"`
	Race       string `json:"race" validate:"oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex        string `json:"sex" validate:"oneof=male female"`
	HasMatched bool   `json:"hasMatched"`
	// ==
	// Age in month example:
	// '=>4': searches data that have more than 4 months
	// '=<4': searches data that have less than 4 months
	// '=4': searches data that have exact 4 month
	AgeInMonth string `json:"ageInMonth"`
	Owned      bool   `json:"owned"`
	Search     string `json:"search"`
}

type CatCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	Race        string `validate:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex         string `validate:"required,oneof=male female"`
	AgeInMonth  int    `validate:"required,numeric,min=1,max=120082"`
	Description string
	ImageUrls   []string `validate:"required,min=1,dive,required,url"`
}

type CatUpdateRequest struct {
	Id          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=30"`
	Race        string `validate:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex         string `validate:"required,oneof=male female"`
	AgeInMonth  int    `validate:"required,numeric,min=1,max=120082"`
	Description string
	ImageUrls   []string `validate:"required,min=1,dive,required,url"`
}
