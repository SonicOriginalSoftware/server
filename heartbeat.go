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

const heartBeatName = "heartbeat"

var (
	gitCommit = ""
)

type heartBeat struct {
	logger logger.Log
}

type data struct {
	Status string `json:"status"`
	Commit string `json:"commit"`
}

func (handler *heartBeat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("%v %v\n", r.Method, r.URL.Path)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		handler.logger.Warn("Unable to read build info\n")
		http.Error(w, "Build info not available", http.StatusInternalServerError)
		return
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			gitCommit = setting.Value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data{Status: "ok", Commit: gitCommit})
}

// RegisterHeartBeat handler
func RegisterHeartBeat(mux *http.ServeMux) (route string) {
	logger := logger.New(
		heartBeatName,
		logger.DefaultSeverity,
		os.Stdout,
		os.Stderr,
	)

	route = fmt.Sprintf("/%s", heartBeatName)
	handler := &heartBeat{logger}
	mux.Handle(route, handler)

	return
}
