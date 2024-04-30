package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	userRepositoryInternal "enigmanations/cats-social/internal/user/repository"
	"enigmanations/cats-social/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
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

func (m *AuthMiddleware) ProtectedHandler(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

		userRepository := userRepositoryInternal.NewUserRepository(m.pool)
		_, err = userRepository.Get(m.ctx, tokenData.Uid)
		if err != nil {
			log.Printf("Invalid token (not found on store): %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Set the HTTP Basic authentication header
		w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", encodedToken))

		// Delegate request to the given handle
		next := h
		next(w, r, ps)
	}
}
