package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bit8bytes/toolbox/logger"
	"github.com/bit8bytes/toolbox/uuid"
)

type Middleware interface {
	LogRequest(next http.Handler) http.Handler
	RecoverPanic(next http.Handler) http.Handler
	AddTraceID(next http.Handler) http.Handler
}

const (
	TraceIDKey = "trace_id"
	SubKey     = "sub"
	NameKey    = "name"
	EmailKey   = "email"
	PictureKey = "picture"
	RolesKey   = "roles"
)

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

func (mw *middleware) AddUserInfoFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, SubKey, "no-x-sub-provided")
		ctx = context.WithValue(ctx, EmailKey, "no-x-email-provided")
		ctx = context.WithValue(ctx, NameKey, "no-x-name-provided")
		ctx = context.WithValue(ctx, PictureKey, "no-x-picture-provided")
		ctx = context.WithValue(ctx, RolesKey, "no-x-roles-provided")

		for key, values := range r.Header {
			if len(values) > 0 {
				switch key {
				case "X-Sub":
					ctx = context.WithValue(ctx, SubKey, values[0])
				case "X-Email":
					ctx = context.WithValue(ctx, EmailKey, values[0])
				case "X-Name":
					ctx = context.WithValue(ctx, NameKey, values[0])
				case "X-Picture":
					ctx = context.WithValue(ctx, PictureKey, values[0])
				case "X-Roles":
					roles := strings.Split(values[0], ",")
					ctx = context.WithValue(ctx, RolesKey, roles)
				}
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
