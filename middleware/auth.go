package middleware

import (
	"log"
	"net/http"

	"enigmanations/cats-social/pkg/jwt"

	"github.com/julienschmidt/httprouter"
)

func ProtectedHandler(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get access token from header
		encodedToken, err := jwt.GetTokenFromAuthHeader(r)
		if err != nil { // error getting Token from auth header
			log.Printf("Error getting token from Header: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Verify access token
		if _, err := jwt.ValidateToken(encodedToken); err != nil {
			log.Printf("Invalid token: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Delegate request to the given handle
		next := h
		next(w, r, ps)
	}
}
