package jwt

import (
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/pkg/env"
	"time"

	"github.com/golang-jwt/jwt"
)

const accessTokenDurationSeconds = 8 * time.Hour // 8 hours

func GenerateAccessToken(userID uint64, credential *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"sub": credential.UUID,
		"exp": time.Now().Add(accessTokenDurationSeconds).Unix(),
	})
	tokenString, err := token.SignedString([]byte(env.GetSecretKey()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateSessionTokenJWT(
	userID uint64,
	credential *user.User,
	session *user.UserSession,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"sub": credential.UUID,
		"sid": session.AccessToken,
		"exp": session.TokenExpires,
	})
	tokenString, err := token.SignedString([]byte(env.GetSecretKey()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
