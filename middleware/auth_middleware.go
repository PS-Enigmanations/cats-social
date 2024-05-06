package middleware

import (
	"log"
	"net/http"

	"enigmanations/cats-social/internal/common/auth"
	userRepository "enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthMiddleware struct {
	pool *pgxpool.Pool
}

func NewAuthMiddleware(pool *pgxpool.Pool) AuthMiddleware {
	return AuthMiddleware{
		pool: pool,
	}
}

func (m *AuthMiddleware) UseAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get access token from header
		encodedToken, err := jwt.GetTokenFromAuthHeader(ctx)
		if err != nil { // error getting Token from auth header
			log.Printf("Error getting token from Header: %v", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Verify access token
		tokenData, err := jwt.ValidateToken(encodedToken)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check user
		userRepository := userRepository.NewUserRepository(m.pool)
		_, err = userRepository.Get(ctx, tokenData.Uid)
		if err != nil {
			log.Printf("Invalid token (not found on store): %v", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Set authentication context
		ctx.Set(auth.AuthorizationPayloadKey, tokenData)

		// Delegate request to the given handle
		ctx.Next()
	}
}
