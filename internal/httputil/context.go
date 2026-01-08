package httputil

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/nadianeyl/nema-api/internal/service"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(r *http.Request, user *service.UserResponse) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetUser(r *http.Request) *service.UserResponse {
	user, ok := r.Context().Value(userContextKey).(*service.UserResponse)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

func ContextGetUserID(r *http.Request) uuid.UUID {
	user := ContextGetUser(r)
	return user.ID
}
