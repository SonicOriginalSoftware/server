//revive:disable:package-comments

package handlers

import (
	"api-server/routes/api"
	"api-server/routes/app"
	"api-server/routes/auth"
	"api-server/routes/git"
	"log"
	"net/http"
)

// SubdomainHandler defines an extension of the http.Handler interface
type SubdomainHandler interface {
	http.Handler
	Prefix() string
	Address() string
}

// NewSubdomainMap returns a mapping of subdomain names and their associated handlers
func NewSubdomainMap(outlog, errlog *log.Logger) []SubdomainHandler {
	apiHandler := api.NewHandler(outlog, errlog)
	appHandler := app.NewHandler(outlog, errlog)
	authHandler := auth.NewHandler(outlog, errlog)
	gitHandler := git.NewHandler(outlog, errlog)

	return []SubdomainHandler{
		apiHandler,
		appHandler,
		authHandler,
		gitHandler,
	}
}
