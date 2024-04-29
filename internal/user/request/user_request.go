package request

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=50"`
}
