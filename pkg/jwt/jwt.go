package jwt

import (
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/pkg/env"
	"time"

	"github.com/golang-jwt/jwt"
)

const AccessTokenDurationSeconds = 8 * time.Hour // 8 hours

func GenerateAccessToken(userID uint64, credential *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userID,
		"sub": credential.Name,
		"exp": time.Now().Add(AccessTokenDurationSeconds).Unix(),
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
		"sub": credential.Name,
		"sid": session.Token,
		"exp": session.ExpiresAt.Second(),
	})
	tokenString, err := token.SignedString([]byte(env.GetSecretKey()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
