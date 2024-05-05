package jwt

import (
	"enigmanations/cats-social/internal/session"
	"enigmanations/cats-social/internal/user"
	"enigmanations/cats-social/pkg/env"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	session *session.Session,
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

type Claims struct {
	ExpiresAt int64 `json:"exp,omitempty"`
	Uid       int   `json:"uid"`
	jwt.MapClaims
}

type TokenData struct {
	ExpiresAt int64 `json:"exp,omitempty"`
	Uid       int   `json:"uid"`
}

func ValidateToken(encodedToken string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(env.GetSecretKey()), nil
	})
	if err != nil {
		return nil, err
	}

	var (
		Uid int
		Exp int64
	)

	// check claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		Uid = claims.Uid
		Exp = claims.ExpiresAt
	} else {
		log.Println("Issue with claims")
		return nil, errors.New("Issue with claims")
	}

	return &TokenData{Uid: Uid, ExpiresAt: Exp}, nil
}

func GetTokenFromAuthHeader(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	// check if Authorization token is set
	if authorizationHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	// Remove bearer in the authorization header
	authorizationHeaderParts := strings.Fields(authorizationHeader)
	if len(authorizationHeaderParts) != 2 || strings.ToLower(authorizationHeaderParts[0]) != "bearer" {
		return "", errors.New("invalid Token - not of type: Bearer")
	}
	return authorizationHeaderParts[1], nil
}
