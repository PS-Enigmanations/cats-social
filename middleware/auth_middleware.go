package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"enigmanations/cats-social/internal/session"
	userRepository "enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthMiddleware struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewAuthMiddleware(pool *pgxpool.Pool, ctx context.Context) AuthMiddleware {
	return AuthMiddleware{
		pool: pool,
		ctx:  ctx,
	}
}

func (m *AuthMiddleware) ProtectedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get access token from header
		encodedToken, err := jwt.GetTokenFromAuthHeader(r)
		if err != nil { // error getting Token from auth header
			log.Printf("Error getting token from Header: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Verify access token
		tokenData, err := jwt.ValidateToken(encodedToken)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userRepository := userRepository.NewUserRepository(m.pool)
		_, err = userRepository.Get(m.ctx, tokenData.Uid)
		if err != nil {
			log.Printf("Invalid token (not found on store): %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Set the HTTP Basic authentication header
		w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", encodedToken))

		// Set authenticaton context
		ctx := session.NewAuthenticationContext(r.Context(), tokenData.Uid)

		// Delegate request to the given handle
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
