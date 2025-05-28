package responder

import (
	"log/slog"
	"net/http"
)

type Envelope map[string]any

type Responder struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Responder {
	return &Responder{
		logger: logger,
	}
}

func (h *Responder) LogError(r *http.Request, err error) {
	var (
		host   = r.Host
		ip     = r.RemoteAddr
		proto  = r.Proto
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	h.logger.Error(
		err.Error(),
		slog.String("host", host),
		slog.String("proto", proto),
		slog.String("ip", ip),
		slog.String("method", method),
		slog.String("uri", uri),
	)
}
