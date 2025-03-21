package responder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bit8bytes/toolbox/logger"
)

type HttpResponder struct {
	logger    logger.Logger
	serviceId string
}

func NewHttp(logger logger.Logger) *HttpResponder {
	return &HttpResponder{
		logger: logger,
	}
}

func (h *HttpResponder) WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (h *HttpResponder) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

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

func (h *HttpResponder) logError(r *http.Request, err error) {
	var (
		service = h.serviceId
		host    = r.Host
		ip      = r.RemoteAddr
		proto   = r.Proto
		method  = r.Method
		uri     = r.URL.RequestURI()
	)

	h.logger.Error(
		err.Error(),
		slog.String("service", service),
		slog.String("host", host),
		slog.String("proto", proto),
		slog.String("ip", ip),
		slog.String("method", method),
		slog.String("uri", uri),
	)
}

func (h *HttpResponder) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}
	err := h.WriteJSON(w, status, env, nil)
	if err != nil {
		h.logError(r, err)
		w.WriteHeader(500)
	}
}

func (h *HttpResponder) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	h.errorResponse(w, r, http.StatusInternalServerError, err.Error())
}

func (h *HttpResponder) NotFound(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	h.errorResponse(w, r, http.StatusNotFound, err.Error())
}

func (h *HttpResponder) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (h *HttpResponder) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	h.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (h *HttpResponder) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *HttpResponder) InvalidBearerAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	// Note: We’re including a WWW-Authenticate: Bearer header here to help inform or remind the client
	// that we expect them to authenticate using a bearer token.
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (h *HttpResponder) InvalidCookieAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	// Note: We’re including a WWW-Authenticate: Bearer header here to help inform or remind the client
	// that we expect them to authenticate using a bearer token.
	w.Header().Set("WWW-Authenticate", "Cookie")

	message := "invalid or missing authentication token"
	h.errorResponse(w, r, http.StatusUnauthorized, message)
}
