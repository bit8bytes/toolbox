package middleware

import (
	"context"
	"net/http"
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

func (m *middleware) GetRoles(r *http.Request) string {
	roles := r.Context().Value(RolesKey).(string)
	if roles == "" {
		roles = NoRolesHeader
	}
	return roles
}
