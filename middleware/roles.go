package middleware

import (
	"context"
	"net/http"
	"strings"
)

const (
	RolesHeader   string = "X-roles"
	NoRolesHeader string = "no-x-roles-povided"
)

func (m *middleware) AddRolesFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roles := r.Header.Get(RolesHeader)
		ctx := context.WithValue(r.Context(), RolesKey, roles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) GetRoles(r *http.Request) []string {
	ctxRoles := r.Context().Value(RolesKey).(string)
	roles := strings.Split(ctxRoles, ",")
	if len(roles) == 0 {
		roles = []string{}
	}
	return roles
}
