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

	RequirePermission(permission string, handler http.HandlerFunc) http.HandlerFunc
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
		chain := final
		for i := len(middlewares) - 1; i >= 0; i-- {
			chain = middlewares[i](chain)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if m.excluded.MatchString(r.URL.Path) {
				final.ServeHTTP(w, r)
			} else {
				chain.ServeHTTP(w, r)
			}
		})
	}
}

// Exclude sets the regular expression to exclude certain paths from the middleware
func (m *middleware) Exclude(excluded *regexp.Regexp) {
	m.excluded = excluded
}
