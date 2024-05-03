package cat

type Cat struct {
	Id          int      `json:"id"`
	UserId      int      `json:"user_id" validate:"required"`
	Name        string   `json:"name" validate:"required,min=1,max=30"`
	Race        Race     `json:"race" validate:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex         Sex      `json:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string   `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" validate:"required,min=1,dive,required,url"`
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

type Race string

const (
	Persian          Race = "Persian"
	MaineCoon        Race = "Maine Coon"
	Siamese          Race = "Siamese"
	Ragdoll          Race = "Ragdoll"
	Bengal           Race = "Bengal"
	Sphynx           Race = "Sphynx"
	BritishShorthair Race = "British Shorthair"
	Abyssinian       Race = "Abyssinian"
	ScottishFold     Race = "Scottish Fold"
	Birman           Race = "Birman"
)

type Sex string

const (
	Male   Sex = "male"
	Female Sex = "female"
)
