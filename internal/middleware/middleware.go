package middleware

import (
	"fmt"
	"net/http"

	"github.com/nadianeyl/nema-api/internal/config"
	"github.com/nadianeyl/nema-api/internal/httputil"
	"github.com/nadianeyl/nema-api/internal/jsonlog"
)

type Middleware struct {
	Config   config.Config
	Logger   *jsonlog.Logger
	HTTPUtil httputil.HTTPUtil
}

func New(cfg config.Config, logger *jsonlog.Logger, httpUtil httputil.HTTPUtil) Middleware {
	return Middleware{
		Config:   cfg,
		Logger:   logger,
		HTTPUtil: httpUtil,
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
