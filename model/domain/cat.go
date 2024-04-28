package domain

type Cat struct {
	Name        string     `json:"name" validate:"required,min=1,max=30"`
	Race        Race       `json:"race" validate:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex         Sex        `json:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int        `json:"ageInMonth" validate:"required,min=1,max=120082"`
	Description string     `json:"description" validate:"required,min=1,max=200"`
	ImageUrls   []ImageUrl `json:"imageUrls" validate:"required,min=1,dive,required,url"`
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

type ImageUrl struct {
	Url string `json:"url" validate:"required,url"`
}
