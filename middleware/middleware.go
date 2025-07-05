package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	TraceIDKey   contextKey = "trace_id"
	UserIDKey    contextKey = "user_id"
	RequestIDKey contextKey = "request_id"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Middleware struct {
	logger         *slog.Logger
	excludedPaths  map[string]bool
	excludedPrefix []string
}

func New(logger *slog.Logger) *Middleware {
	if logger == nil {
		logger = slog.Default()
	}

	mw := &Middleware{
		logger:         logger,
		excludedPaths:  make(map[string]bool),
		excludedPrefix: make([]string, 0),
	}

	return mw
}

// Chain combines middlewares and returns a handler.
// Build chain from right to left
func (m *Middleware) Chain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		// Build chain from right to left
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}

// ExcludePaths excludes exact path matches
// Example: mw.ExcludePaths("/health")
func (m *Middleware) ExcludePaths(paths ...string) {
	for _, path := range paths {
		m.excludedPaths[path] = true
	}
}

// ExcludePrefixes excludes paths starting with given prefixes
// Example: mw.ExcludePrefixes("/health/")
func (m *Middleware) ExcludePrefixes(prefixes ...string) {
	m.excludedPrefix = append(m.excludedPrefix, prefixes...)
}

// shouldSkip checks if request should skip middleware
func (m *Middleware) ShouldSkip(r *http.Request) bool {
	path := r.URL.Path

	// Check exact path matches (O(1) lookup)
	if m.excludedPaths[path] {
		return true
	}

	// Check prefix matches (O(n) but typically very small n)
	for _, prefix := range m.excludedPrefix {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}

// LogRequest logs HTTP requests with response details
func (m *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.ShouldSkip(r) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		traceID := getTraceID(r)

		// Add trace ID to context
		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status and size
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		m.logger.Info("request",
			slog.String("trace_id", traceID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Int("size", wrapped.size),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// RecoverPanic recovers from panics
func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				traceID := GetTraceIDFromContext(r.Context())

				m.logger.Error("panic recovered",
					slog.String("trace_id", traceID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("error", err),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"internal server error","trace_id":"%s"}`, traceID)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// responseWriter captures response metadata
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

// GetTraceIDFromContext returns the trace id from the context
func GetTraceIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return "unknown"
}

// Helper functions
func getTraceID(r *http.Request) string {
	headers := []string{"X-Trace-Id", "X-Request-Id", "X-Correlation-Id"}

	for _, header := range headers {
		if id := r.Header.Get(header); id != "" {
			return id
		}
	}

	return fmt.Sprintf("%d", time.Now().UnixNano())
}
