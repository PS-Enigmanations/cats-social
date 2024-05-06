package auth

import (
	"enigmanations/cats-social/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// GetCurrentUser is a helper function to get the auth payload from the context
func GetCurrentUser(ctx *gin.Context) *jwt.TokenData {
	return ctx.MustGet(AuthorizationPayloadKey).(*jwt.TokenData)
}
