// Package gzip provides gzip compression middleware
package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/bit8bytes/toolbox/middleware"
)

// Config holds gzip middleware configuration
type Config struct {
	// Compression level (1-9, where 9 is best compression)
	Level int
	// Minimum response size to compress (bytes)
	MinSize int
	// Content types to compress (if empty, compresses all)
	Types []string
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Level:   gzip.DefaultCompression,
		MinSize: 1024, // 1KB
		Types: []string{
			"text/html",
			"text/css",
			"text/javascript",
			"application/javascript",
			"application/json",
			"application/xml",
			"text/xml",
			"text/plain",
		},
	}
}

// New creates gzip middleware with custom config
func New(mw *middleware.Middleware, config *Config) middleware.MiddlewareFunc {
	if config == nil {
		config = DefaultConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if excluded or client doesn't accept gzip
			if mw.ShouldSkip(r) || !acceptsGzip(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Create gzip writer
			gz, err := gzip.NewWriterLevel(w, config.Level)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			defer gz.Close()

			// Set headers
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")

			// Wrap response writer
			wrapped := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gz,
				config:         config,
			}

			next.ServeHTTP(wrapped, r)
		})
	}
}

// Handler creates gzip middleware with default config
func Handler(mw *middleware.Middleware) middleware.MiddlewareFunc {
	return New(mw, DefaultConfig())
}

// gzipResponseWriter wraps response writer for gzip compression
type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
	config         *Config
	headerWritten  bool
	shouldCompress bool
}

func (grw *gzipResponseWriter) WriteHeader(code int) {
	if grw.headerWritten {
		return
	}
	grw.headerWritten = true

	// Check if we should compress based on content type
	if len(grw.config.Types) > 0 {
		contentType := grw.Header().Get("Content-Type")
		grw.shouldCompress = shouldCompressType(contentType, grw.config.Types)
	} else {
		grw.shouldCompress = true
	}

	grw.ResponseWriter.WriteHeader(code)
}

func (grw *gzipResponseWriter) Write(data []byte) (int, error) {
	if !grw.headerWritten {
		grw.WriteHeader(http.StatusOK)
	}

	// Check minimum size
	if len(data) < grw.config.MinSize && !grw.shouldCompress {
		return grw.ResponseWriter.Write(data)
	}

	if grw.shouldCompress {
		return grw.Writer.Write(data)
	}
	return grw.ResponseWriter.Write(data)
}

func acceptsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func shouldCompressType(contentType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if strings.HasPrefix(contentType, allowed) {
			return true
		}
	}
	return false
}
