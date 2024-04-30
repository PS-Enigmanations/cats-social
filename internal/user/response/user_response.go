package response

import "enigmanations/cats-social/internal/user"

type UserShow struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserCreateResponse struct {
	Message string   `json:"message"`
	Data    UserShow `json:"data"`
}

const UserCreateSuccMessage = "User registered successfully"

func UserToUserCreateResponse(u user.User) *UserCreateResponse {
	return &UserCreateResponse{
		Message: UserCreateSuccMessage,
		Data: UserShow{
			Email: u.Email,
			Name:  u.Name,
		},
	}
}
