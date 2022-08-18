//revive:disable:package-comments

package main

import (
	"api-server/lib"
	"api-server/lib/net/handlers"
	"api-server/routes/api"
	"api-server/routes/app"
	"api-server/routes/auth"
	"api-server/routes/git"
	"os"
)

func main() {
	apiHandler := api.NewHandler()
	appHandler := app.NewHandler()
	authHandler := auth.NewHandler()
	gitHandler := git.NewHandler()

	subdomains := []handlers.SubdomainHandler{
		apiHandler,
		appHandler,
		authHandler,
		gitHandler,
	}

	defer func() { os.Exit(lib.Run(subdomains)) }()
}
