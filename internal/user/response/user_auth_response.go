package response

import "enigmanations/cats-social/internal/user"

type UserLoginResponse struct {
	Message string   `json:"message"`
	Data    UserShow `json:"data"`
}

const UserLoginSuccMessage = "User logged successfully"

func UserToUserLoginResponse(u user.User, uSession user.UserSession) *UserLoginResponse {
	return &UserLoginResponse{
		Message: UserLoginSuccMessage,
		Data: UserShow{
			Email:       u.Email,
			Name:        u.Name,
			AccessToken: uSession.Token,
		},
	}
}
