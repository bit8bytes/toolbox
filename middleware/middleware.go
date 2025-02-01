package middleware

import (
	"net/http"

	"github.com/bit8bytes/toolbox/logger"
)

type Middleware interface {
	Chain(middlewares ...middlewares) middlewares
	LogRequest(next http.Handler) http.Handler
	RecoverPanic(next http.Handler) http.Handler
	AddTraceIdFromHeaderToContext(next http.Handler) http.Handler
	AddUserInfoFromHeaderToContext(next http.Handler) http.Handler
	AddOrgIdFromHeaderToContext(next http.Handler) http.Handler
}

type middlewares func(http.Handler) http.Handler

type middleware struct {
	logger logger.Logger
}

func NewMiddleware(l logger.Logger) *middleware {
	return &middleware{
		logger: l,
	}
}

func (m *middleware) Chain(middlewares ...middlewares) middlewares {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
