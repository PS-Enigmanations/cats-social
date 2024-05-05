package response

import "enigmanations/cats-social/internal/user"

type SessionLoginShow struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"accessToken"`
}
type SessionLoginResponse struct {
	Message string           `json:"message"`
	Data    SessionLoginShow `json:"data"`
}

const SessionLoginSuccMessage = "User logged successfully"

func SessionToLoginResponse(u user.User, token string) *SessionLoginResponse {
	return &SessionLoginResponse{
		Message: SessionLoginSuccMessage,
		Data: SessionLoginShow{
			Email: u.Email,
			Name:  u.Name,
			Token: token,
		},
	}
}
