package middleware

import (
	"context"
	"net/http"
)

func (m *middleware) AddTraceIdFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := extractTraceIdFromXRequestHeader(r)

		ctx := context.WithValue(r.Context(), TraceIdKey, traceId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTraceIdFromXRequestHeader(r *http.Request) string {
	traceId := "no-x-request-id-header-provided"
	for _, value := range r.Header["X-Request-Id"] {
		traceId = value
		break
	}
	return traceId
}
