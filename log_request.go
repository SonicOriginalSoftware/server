package server

import (
	"context"
	"log/slog"
	"net/http"

	"git.sonicoriginal.software/server/v2/logging"
)

var (
	requestLogger = logging.JSONLogger.With(slog.String("handler", "request"))
)

// LogRequest logs the details of an HTTP request
func LogRequest(ctx context.Context, r *http.Request, logger logging.SLogger) {
	remoteAddr := r.RemoteAddr
	origin := r.Header.Get("Origin")
	agent := r.Header.Get("User-Agent")
	query := r.URL.Query()

	if logger == nil {
		logger = requestLogger
	}

	logger.InfoContext(
		ctx,
		"Request",
		slog.String("remoteAddr", remoteAddr),
		slog.String("origin", origin),
		slog.String("agent", agent),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("query", query.Encode()),
	)
}
