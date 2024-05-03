package request

type CatGetAllQueryParams struct {
	Id         string `form:"id"`
	Limit      string `form:"limit" default:"5"`
	Offset     string `form:"offset" default:"0"`
	Race       string `form:"race" validate:"oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex        string `form:"sex" validate:"oneof=male female"`
	HasMatched string `form:"hasMatched"`
	// ==
	// Age in month example:
	// '=>4': searches data that have more than 4 months
	// '=<4': searches data that have less than 4 months
	// '=4': searches data that have exact 4 month
	AgeInMonth string `json:"ageInMonth"`
	Owned      string `json:"owned"`
	Search     string `json:"search"`
}

type CatCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	Race        CatRace  `validate:"required"`
	Sex         string `validate:"required,oneof=male female"`
	AgeInMonth  int    `validate:"required,numeric,min=1,max=120082"`
	Description string `validate:"required"`
	ImageUrls   []string `validate:"required,min=1,dive,required,url"`
}

type CatRace string

const (
    Persian          CatRace = "Persian"
    MaineCoon        CatRace = "Maine Coon"
    Siamese          CatRace = "Siamese"
    Ragdoll          CatRace = "Ragdoll"
    Bengal           CatRace = "Bengal"
    Sphynx           CatRace = "Sphynx"
    BritishShorthair CatRace = "British Shorthair"
    Abyssinian       CatRace = "Abyssinian"
    ScottishFold     CatRace = "Scottish Fold"
    Birman           CatRace = "Birman"
)

type CatUpdateRequest struct {
	Id          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=30"`
	Race        string `validate:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex         string `validate:"required,oneof=male female"`
	AgeInMonth  int    `validate:"required,numeric,min=1,max=120082"`
	Description string
	ImageUrls   []string `validate:"required,min=1,dive,required,url"`
}
