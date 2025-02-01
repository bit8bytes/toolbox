package middleware

import (
	"context"
	"net/http"
)

func (m *middleware) AddTraceIdFromHeaderToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := getTraceIdFromHeader(r)
		ctx := context.WithValue(r.Context(), TraceIdKey, traceId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) GetTraceIdFromContext(r *http.Request) string {
	traceId := r.Context().Value(TraceIdKey).(string)
	if traceId == "" {
		traceId = "no-x-request-id"
	}
	return traceId
}

func getTraceIdFromHeader(r *http.Request) string {
	return r.Header.Get("X-Request-Id")
}
