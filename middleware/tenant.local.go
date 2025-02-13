//go:build local

package middleware

import (
	"context"
	"net/http"
)

const (
	TenantIdHeader               string = "X-tenant-id"
	TenantDisplayNameHeader      string = "X-tenant-display-name"
	ErrTenantIdRequired          string = "Tenant ID is required"
	ErrTenantDisplayNameRequired string = "Tenant display name is required"
)

func (middleware *middleware) AddTenantIdFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleware.excluded.MatchString(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		tenantId := r.Header.Get(TenantIdHeader)
		if tenantId == "" {
			tenantId = "no-x-tenant-id-local"
		}

		ctx := context.WithValue(r.Context(), TenantIdKey, tenantId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *middleware) GetTenantIdFromContext(r *http.Request) string {
	return r.Context().Value(TenantIdKey).(string)
}

func (middleware *middleware) AddTenantDisplayNameFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleware.excluded.MatchString(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		tenantDisplayName := r.Header.Get(TenantDisplayNameHeader)
		if tenantDisplayName == "" {
			tenantDisplayName = "no-x-tenant-display-name-local"
			return
		}

		ctx := context.WithValue(r.Context(), TenantDisplayNameKey, tenantDisplayName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *middleware) GetTenantDisplayNameFromContext(r *http.Request) string {
	return r.Context().Value(TenantDisplayNameKey).(string)
}
