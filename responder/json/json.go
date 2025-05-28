// Package json provides utilities for handling JSON requests and responses in HTTP handlers.
// It offers structured error handling for JSON parsing and marshaling with configurable options.
package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bit8bytes/toolbox/responder"
)

const (
	DefaultMaxBytes = 1_048_576 // Default max request body size (1MB)
)

// JSONResponder handles JSON encoding and decoding for HTTP requests and responses.
// It provides structured error handling and configurable request body limits.
type JSONResponder struct {
	logger   *slog.Logger
	maxBytes int64
	responder.Responder
}

// Options is a function type for configuring JSONResponder instances.
type Options func(*JSONResponder)

// WithMaxBytes returns an option that sets the maximum request body size in bytes.
// This limit is enforced when reading JSON request bodies to prevent excessive memory usage.
func WithMaxBytes(b int64) Options {
	return func(jr *JSONResponder) {
		jr.maxBytes = b
	}
}

// New creates a new JSONResponder with the provided logger and optional configuration.
// The logger is used for structured logging throughout the JSON handling process.
// If no options are provided, DefaultMaxBytes (1MB) is used as the request body limit.
func New(logger *slog.Logger, opts ...Options) *JSONResponder {
	jr := &JSONResponder{
		logger:   logger,
		maxBytes: DefaultMaxBytes,
	}

	for _, opt := range opts {
		opt(jr)
	}

	return jr
}

// WriteJSON encodes the provided data as JSON and writes it to the HTTP response.
// It sets the Content-Type header to application/json and applies any custom headers.
// The data parameter should be of type responder.Envelope for consistent response structure.
func (jr *JSONResponder) WriteJSON(w http.ResponseWriter, status int, data responder.Envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Apply custom headers first
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, writeErr := w.Write(js)

	return writeErr
}

// ReadJSON reads and decodes JSON from the HTTP request body into the provided destination.
// It enforces the configured maximum body size and provides detailed error messages
// for various JSON parsing failures including syntax errors, type mismatches, and unknown fields.
// The method ensures only a single JSON value is present in the request body.
func (jr *JSONResponder) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, jr.maxBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// ServerErrorResponse sends a 500 Internal Server Error response with the error message.
// It logs the error details and returns a JSON error response to the client.
func (jr *JSONResponder) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	jr.LogError(r, err)
	jr.errorResponse(w, r, http.StatusInternalServerError, err.Error())
}

// NotFound sends a 404 Not Found response with the error message.
// It logs the error details and returns a JSON error response to the client.
func (jr *JSONResponder) NotFound(w http.ResponseWriter, r *http.Request, err error) {
	jr.LogError(r, err)
	jr.errorResponse(w, r, http.StatusNotFound, err.Error())
}

// BadRequestResponse sends a 400 Bad Request response with the error message.
// It returns a JSON error response to the client without logging the error.
func (jr *JSONResponder) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	jr.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// FailedValidationResponse sends a 422 Unprocessable Entity response with validation errors.
// The errors parameter should contain field names mapped to their validation error messages.
func (jr *JSONResponder) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	jr.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// InvalidCredentialsResponse sends a 401 Unauthorized response for invalid login credentials.
// This is typically used when username/password authentication fails.
func (jr *JSONResponder) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	jr.errorResponse(w, r, http.StatusUnauthorized, message)
}

// InvalidBearerAuthenticationTokenResponse sends a 401 Unauthorized response for invalid bearer tokens.
// It sets the WWW-Authenticate: Bearer header to inform the client that bearer token authentication is required.
func (jr *JSONResponder) InvalidBearerAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	// Note: We're including a WWW-Authenticate: Bearer header here to help inform or remind the client
	// that we expect them to authenticate using a bearer token.
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	jr.errorResponse(w, r, http.StatusUnauthorized, message)
}

// InvalidCookieAuthenticationTokenResponse sends a 401 Unauthorized response for invalid cookie tokens.
// It sets the WWW-Authenticate: Cookie header to inform the client that cookie authentication is required.
func (jr *JSONResponder) InvalidCookieAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	// Note: We're including a WWW-Authenticate: Cookie header here to help inform or remind the client
	// that we expect them to authenticate using a cookie token.
	w.Header().Set("WWW-Authenticate", "Cookie")
	message := "invalid or missing authentication token"
	jr.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (jr *JSONResponder) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := responder.Envelope{"error": message}
	err := jr.WriteJSON(w, status, env, nil)
	if err != nil {
		jr.LogError(r, err)
		w.WriteHeader(500)
	}
}
