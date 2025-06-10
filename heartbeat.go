//revive:disable:package-comments
package server

import (
	"encoding/json"
	"net/http"

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
