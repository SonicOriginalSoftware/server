//revive:disable:package-comments
package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type heartBeat struct {
	logger    *slog.Logger
	errLogger *slog.Logger
	commit    string
}

type data struct {
	Status string `json:"status"`
	Commit string `json:"commit"`
}

func (handler *heartBeat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LogRequest(r, handler.logger)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data{Status: "ok", Commit: handler.commit}); err != nil {
		handler.errLogger.Error(fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterHeartBeat handler
func RegisterHeartBeat(mux *http.ServeMux, parentJSONLogger, parentJSONErrorLogger *slog.Logger) {
	if parentJSONLogger == nil {
		parentJSONLogger = JSONLogger
	}
	if parentJSONErrorLogger == nil {
		parentJSONErrorLogger = JSONErrorLogger
	}

	logger := parentJSONLogger.With(slog.String("handler", HeartBeatName))
	errorLogger := parentJSONErrorLogger.With(slog.String("handler", HeartBeatName))

	commit := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			TextLogger.Debug(
				"Build info",
				slog.String("key", setting.Key),
				slog.String("value", setting.Value),
			)
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
			}
		}
	} else {
		TextLogger.Warn("Unable to read build info")
	}

	route := fmt.Sprintf("/%s", HeartBeatName)
	handler := &heartBeat{logger, errorLogger, commit}
	mux.Handle(route, handler)

	RegisterLogger.Info("route", slog.String("path", route))
}
