package response

import "enigmanations/cats-social/internal/user"

type UserLoginShow struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"accessToken"`
}
type UserLoginResponse struct {
	Message string        `json:"message"`
	Data    UserLoginShow `json:"data"`
}

const UserLoginSuccMessage = "User logged successfully"

func UserToUserLoginResponse(u user.User, token string) *UserLoginResponse {
	return &UserLoginResponse{
		Message: UserLoginSuccMessage,
		Data: UserLoginShow{
			Email: u.Email,
			Name:  u.Name,
			Token: token,
		},
	}
}
