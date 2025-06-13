//revive:disable:package-comments
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	server_context "git.sonicoriginal.software/server/v2/context"
	"git.sonicoriginal.software/server/v2/logging"
)

var (
	heartBeatLogger logging.SLogger = logging.JSONLogger.With(slog.String("handler", HeartBeatName))
)

type heartBeat struct {
	commit string
}

type data struct {
	Status string `json:"status"`
	Commit string `json:"commit"`
}

func (handler *heartBeat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = server_context.WithUUID(ctx)
	LogRequest(ctx, r, nil)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data{Status: "ok", Commit: handler.commit}); err != nil {
		heartBeatLogger.ErrorContext(ctx, fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RegisterHeartBeat handler
func RegisterHeartBeat(ctx context.Context, mux *http.ServeMux, logger logging.SLogger) {
	ctx = server_context.WithID(ctx, HeartBeatName)
	if logger == nil {
		heartBeatLogger = logger
	}

	commit := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			logging.RegisterLogger.DebugContext(
				ctx,
				"Build info",
				slog.String("key", setting.Key),
				slog.String("value", setting.Value),
			)

			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}
		}
	} else {
		logging.RegisterLogger.WarnContext(ctx, "Unable to read build info")
	}

	route := fmt.Sprintf("/%s", HeartBeatName)
	handler := &heartBeat{commit}
	mux.Handle(route, handler)

	logging.RegisterLogger.InfoContext(ctx, "Registered", slog.String("path", route))
}
