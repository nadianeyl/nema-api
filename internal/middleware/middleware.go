package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
	"github.com/nadianeyl/nema-api/internal/repository"
	"github.com/nadianeyl/nema-api/internal/service"
	"github.com/nadianeyl/nema-api/internal/validator"
)

type Middleware struct {
	Config      config.Config
	Logger      *jsonlog.Logger
	HTTPUtil    httputil.HTTPUtil
	UserService service.UserService
}

func New(cfg config.Config, logger *jsonlog.Logger, httpUtil httputil.HTTPUtil, userService service.UserService) Middleware {
	return Middleware{
		Config:      cfg,
		Logger:      logger,
		HTTPUtil:    httpUtil,
		UserService: userService,
	}
}

func (m *Middleware) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Logger.LogInfo("request", map[string]string{
			"network_address": r.RemoteAddr,
			"protocol":        r.Proto,
			"request_method":  r.Method,
			"request_url":     r.URL.RequestURI(),
		})

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range m.Config.Cors.TrustedOrigins {
				if origin == m.Config.Cors.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						w.Header().Set("Access-Control-Max-Age", "10")

						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.HTTPUtil.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = contextSetUser(r, service.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.HTTPUtil.InvalidAuthTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()
		if service.ValidateTokenPlaintext(v, token); !v.Valid() {
			m.HTTPUtil.InvalidAuthTokenResponse(w, r)
			return
		}

		user, err := m.UserService.GetForToken(repository.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				m.HTTPUtil.InvalidAuthTokenResponse(w, r)
			default:
				m.HTTPUtil.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contextGetUser(r)

		if user.IsAnonymous() {
			m.HTTPUtil.AuthenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contextGetUser(r)

		if !user.Activated {
			m.HTTPUtil.InactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return m.RequireAuthenticatedUser(fn)
}
