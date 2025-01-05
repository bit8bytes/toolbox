package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bit8bytes/toolbox/logger"
	"github.com/bit8bytes/toolbox/uuid"
)

type Middleware interface {
	LogRequest(next http.Handler) http.Handler
	RecoverPanic(next http.Handler) http.Handler
	AddTraceID(next http.Handler) http.Handler
}

const TraceIDKey = "trace_id"

type middlewares func(http.Handler) http.Handler

type middleware struct {
	logger logger.Logger
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

func (m *middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := extractTraceIDFromXRequestHeader(r)

		m.logger.Info("received request",
			slog.String("trace_id", traceID),
			slog.String("host", r.Host),
			slog.String("proto", r.Proto),
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()),
		)
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, fmt.Sprintf("panic: %v", err), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) AddTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, err := uuid.New()
		if err != nil {
			// fallback to prevent middleware from failing
			traceID = uuid.UUID([]byte("unkown"))
		}
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *middleware) AddTraceIDFromXRequestIdHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := extractTraceIDFromXRequestHeader(r)

		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTraceIDFromXRequestHeader(r *http.Request) string {
	traceID := "no-x-request-id-header-provided"
	for _, value := range r.Header["X-Request-Id"] {
		traceID = value
		break
	}
	return traceID
}
