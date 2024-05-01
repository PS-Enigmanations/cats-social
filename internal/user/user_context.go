package user

import (
	"context"
	"net/http"
)

type ctxKey int

type userValue struct {
	Uid int
}

const (
	CtxKeyAuth ctxKey = iota
)

// NewAuthenticationContext initializes the context for protected requests
func NewAuthenticationContext(ctx context.Context, uid int) context.Context {
	return context.WithValue(ctx, CtxKeyAuth, &userValue{
		Uid: uid,
	})
}

// IsAuthenticated checks if this request was authenticated via a middleware
func IsAuthenticated(r *http.Request) bool {
	if v := r.Context().Value(CtxKeyAuth); v != nil {
		return true
	}
	if v := r.Context().Value(CtxKeyAuth); v != nil {
		return true
	}
	return false
}

// GetCurrentUser is a wrapper over the private User context key
func GetCurrentUser(ctx context.Context) *userValue {
	val := ctx.Value(CtxKeyAuth)
	v, ok := val.(*userValue)
	if !ok {
		return nil
	}
	return v
}
