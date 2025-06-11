//revive:disable:package-comments
package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
)

type heartBeat struct {
	logger *slog.Logger
	commit string
}

type data struct {
	Status string `json:"status"`
	Commit string `json:"commit"`
}

func (handler *heartBeat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LogRequest(r, handler.logger, r.URL.Query())

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data{Status: "ok", Commit: handler.commit}); err != nil {
		handler.logger.Error(fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterHeartBeat handler
func RegisterHeartBeat(mux *http.ServeMux, beat *slog.Logger) {
	if beat == nil {
		opts := &slog.HandlerOptions{}
		beat = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	commit := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
			}
		}
	} else {
		beat.Warn("Unable to read build info")
	}

	route := fmt.Sprintf("/%s", HeartBeatName)
	handler := &heartBeat{beat, commit}
	mux.Handle(route, handler)

	beat.Info(fmt.Sprintf("Handler registered for route: %s", route))

	return
}
