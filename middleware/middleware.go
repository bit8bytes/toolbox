package middleware

import (
	"net/http"
	"regexp"

	"github.com/bit8bytes/toolbox/logger"
)

type Middleware interface {
	Chain(middlewares ...middlewares) middlewares
	Exclude(excluded *regexp.Regexp)
	LogRequest(next http.Handler) http.Handler
	RecoverPanic(next http.Handler) http.Handler

	AddTraceIdFromHeaderToContext(next http.Handler) http.Handler
	GetTraceIdFromContext(r *http.Request) string

	AddOrgIdForUserFromHeaderToContext(next http.Handler) http.Handler
	GetOrgIdForUserFromContext(r *http.Request) string
	AddUserSubFromHeaderToContext(next http.Handler) http.Handler
	GetUserSubFromContext(r *http.Request) string
	AddUserNameFromHeaderToContext(next http.Handler) http.Handler
	GetUserNameFromContext(r *http.Request) string
	AddUserNicknameFromHeaderToContext(next http.Handler) http.Handler
	GetUserNicknameFromContext(r *http.Request) string
	AddUserEmailFromHeaderToContext(next http.Handler) http.Handler
	GetUserEmailFromContext(r *http.Request) string
	AddUserPictureFromHeaderToContext(next http.Handler) http.Handler
	GetUserPictureFromContext(r *http.Request) string

	AddTenantIdFromHeaderToContext(next http.Handler) http.Handler
	GetTenantIdFromContext(r *http.Request) string
	AddTenantDisplayNameFromHeaderToContext(next http.Handler) http.Handler
	GetTenantDisplayNameFromContext(r *http.Request) string

	AddRolesFromHeaderToContext(next http.Handler) http.Handler
	GetRoles(r *http.Request) []string

	Gzip(next http.Handler) http.Handler
}

type middlewares func(http.Handler) http.Handler

type middleware struct {
	logger   logger.Logger
	excluded *regexp.Regexp
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

// Exclude sets the regular expression to exclude certain paths from the middleware
func (m *middleware) Exclude(excluded *regexp.Regexp) {
	m.excluded = excluded
}
