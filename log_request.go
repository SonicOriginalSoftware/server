package server

import (
	"log/slog"
	"net/http"
	"net/url"
)

// LogRequest logs the details of an HTTP request
func LogRequest(r *http.Request, logger *slog.Logger, query url.Values) {
	remoteAddr := r.RemoteAddr
	origin := r.Header.Get("Origin")
	agent := r.Header.Get("User-Agent")
	logger.Info(
		remoteAddr,
		"origin", origin,
		"agent", agent,
		"method", r.Method,
		"path", r.URL.Path,
		"query", query.Encode(),
	)
}
