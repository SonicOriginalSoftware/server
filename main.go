//revive:disable:package-comments

package main

import (
	"os"
	"server/lib"
	"server/lib/handlers"
	"server/routes/api"
	"server/routes/app"
	"server/routes/auth"
	"server/routes/git"
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
