package middleware

import (
	"context"
	"net/http"
	"strings"
)

// Available user info: sub, name, nickname, email, email verified, picture, roles.
// Usage: i.e. middleware.SubKey
func (mw *middleware) AddUserInfoFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, SubKey, "no-x-sub-provided")
		ctx = context.WithValue(ctx, NameKey, "no-x-name-provided")
		ctx = context.WithValue(ctx, NicknameKey, "no-x-nickname-provided")
		ctx = context.WithValue(ctx, EmailKey, "no-x-email-provided")
		ctx = context.WithValue(ctx, EmailVerifiedKey, "no-x-email-verified-provided")
		ctx = context.WithValue(ctx, PictureKey, "no-x-picture-provided")
		ctx = context.WithValue(ctx, RolesKey, "no-x-roles-provided")

		for key, values := range r.Header {
			if len(values) > 0 {
				switch key {
				case "X-sub":
					ctx = context.WithValue(ctx, SubKey, values[0])
				case "X-name":
					ctx = context.WithValue(ctx, NameKey, values[0])
				case "X-nickname":
					ctx = context.WithValue(ctx, NicknameKey, values[0])
				case "X-email":
					ctx = context.WithValue(ctx, EmailKey, values[0])
				case "X-email-verified":
					ctx = context.WithValue(ctx, EmailVerifiedKey, values[0])
				case "X-picture":
					ctx = context.WithValue(ctx, PictureKey, values[0])
				case "X-roles":
					roles := strings.Split(values[0], ",")
					ctx = context.WithValue(ctx, RolesKey, roles)
				}
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *middleware) AddOrgIdFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgId := extractOrgIdFromXHeader(r)

		ctx := context.WithValue(r.Context(), OrgIdKey, orgId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractOrgIdFromXHeader(r *http.Request) string {
	orgId := "no-x-request-org-id-header-provided"
	for _, value := range r.Header["X-org-id"] {
		orgId = value
		break
	}
	return orgId
}
