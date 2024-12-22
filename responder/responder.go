package responder

import "net/http"

type Envelope map[string]any

type Responder interface {
	ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error
	WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error

	ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
	FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string)
	InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request)
	InvalidBearerAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request)
	InvalidCookieAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request)

	logError(r *http.Request, err error)
	errorResponse(w http.ResponseWriter, r *http.Request, status int, message any)
}
