//revive:disable:package-comments
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"git.sonicoriginal.software/logger/v2"
)

type heartBeat struct {
	logger logger.Log
	commit string
}

type data struct {
	Status string `json:"status"`
	Commit string `json:"commit"`
}

func (handler *heartBeat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("%v %v\n", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data{Status: "ok", Commit: handler.commit})
}

// RegisterHeartBeat handler
func RegisterHeartBeat(mux *http.ServeMux) {
	commit := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
			}
		}
	} else {
		logger.DefaultLogger.Warn("Unable to read build info\n")
	}

	beat := logger.New(
		HeartBeatName,
		logger.DefaultSeverity,
		os.Stdout,
		os.Stderr,
	)

	route := fmt.Sprintf("/%s", HeartBeatName)
	handler := &heartBeat{beat, commit}
	mux.Handle(route, handler)

	beat.Info("Handler registered for route [%v]\n", route)

	return
}
