// Package cors provides CORS (Cross-Origin Resource Sharing) middleware
package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bit8bytes/toolbox/middleware"
)

// Config holds CORS configuration
type Config struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	ExposedHeaders     []string
	AllowCredentials   bool
	MaxAge             int // in seconds
	OptionsPassthrough bool
}

// DefaultConfig returns sensible CORS defaults
func DefaultConfig() *Config {
	return &Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:     []string{},
		AllowCredentials:   false,
		MaxAge:             86400, // 24 hours
		OptionsPassthrough: false,
	}
}

// New creates CORS middleware with custom config
func New(mw *middleware.Middleware, config *Config) middleware.MiddlewareFunc {
	if config == nil {
		config = DefaultConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if excluded
			if mw.ShouldSkip(r) {
				next.ServeHTTP(w, r)
				return
			}

			origin := r.Header.Get("Origin")

			// Set CORS headers if origin is allowed
			if isOriginAllowed(origin, config.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)

				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				if len(config.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
				}
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
				}

				if !config.OptionsPassthrough {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Handler creates CORS middleware with default config
func Handler(mw *middleware.Middleware) middleware.MiddlewareFunc {
	return New(mw, DefaultConfig())
}

// AllowAll creates permissive CORS middleware (development only!)
func AllowAll(mw *middleware.Middleware) middleware.MiddlewareFunc {
	config := &Config{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"*"},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   false, // Can't be true with wildcard origins
		MaxAge:             86400,
		OptionsPassthrough: false,
	}
	return New(mw, config)
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
