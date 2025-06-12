package server

import (
	"log/slog"
	"net/http"
)

// LogRequest logs the details of an HTTP request
func LogRequest(r *http.Request, logger *slog.Logger) {
	remoteAddr := r.RemoteAddr
	origin := r.Header.Get("Origin")
	agent := r.Header.Get("User-Agent")
	query := r.URL.Query()

	if logger == nil {
		logger = JSONLogger
	}

	logger.Info(
		"Request",
		slog.String("remoteAddr", remoteAddr),
		slog.String("origin", origin),
		slog.String("agent", agent),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("query", query.Encode()),
	)
}
