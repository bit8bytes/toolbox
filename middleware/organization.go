package middleware

import (
	"context"
	"net/http"
)

func (middleware *middleware) AddOrgIdForUserFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgId := r.Header.Get("X-Org-Id")
		if orgId == "" {
			orgId = "no-x-org-id"
		}

		ctx := context.WithValue(r.Context(), OrgIdKey, orgId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (middleware *middleware) GetOrgIdForUserFromContext(r *http.Request) string {
	orgId := r.Context().Value(OrgIdKey).(string)
	if orgId == "" {
		orgId = "no-x-org-id"
	}
	return orgId
}
