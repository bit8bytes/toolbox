package middleware

import (
	"context"
	"net/http"
)

const (
	TenantIdHeader               = "X-Tenant-Id"
	TenantDisplayNameHeader      = "X-Tenant-Display-Name"
	ErrTenantIdRequired          = "Tenant ID is required"
	ErrTenantDisplayNameRequired = "Tenant display name is required"
)

func (middleware *middleware) AddTenantIdForUserFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantId := r.Header.Get(TenantIdHeader)
		if tenantId == "" {
			http.Error(w, ErrTenantIdRequired, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), TenantIdKey, tenantId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *middleware) GetTenantIdForUserFromContext(r *http.Request) string {
	return r.Context().Value(TenantIdKey).(string)
}

func (middleware *middleware) AddTenantDisplayNameFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantDisplayName := r.Header.Get(TenantDisplayNameHeader)
		if tenantDisplayName == "" {
			http.Error(w, ErrTenantDisplayNameRequired, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), TenantDisplayNameKey, tenantDisplayName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *middleware) GetTenantTenantNameFromContext(r *http.Request) string {
	return r.Context().Value(TenantDisplayNameKey).(string)
}
