package middleware

import (
	"context"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/service"
)

type contextKey string

const userContextKey = contextKey("user")

func contextSetUser(r *http.Request, user *service.UserResponse) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func contextGetUser(r *http.Request) *service.UserResponse {
	user, ok := r.Context().Value(userContextKey).(*service.UserResponse)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
