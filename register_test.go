//revive:disable:package-comments

package server_test

import (
	"net/http"
	"testing"

	"git.sonicoriginal.software/logger"
	"git.sonicoriginal.software/server/v2"
)

type handler struct {
	http.Handler
}

func TestRegisterRootHandler(t *testing.T) {
	path := "/"
	server.RegisterHandler(path, handler{})
	logger.DefaultLogger.Info("service registered for route: %v", path)
}

func TestRegisterServiceHandler(t *testing.T) {
	path := "service"
	route := server.RegisterHandler(path, handler{})
	logger.DefaultLogger.Info("service registered for route: %v", route)
}
