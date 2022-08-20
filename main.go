//revive:disable:package-comments

package main

import (
	"os"
	"server/handlers"
	"server/lib"
	"server/routes/api"
	"server/routes/app"
	"server/routes/auth"
	"server/routes/git"
	"server/routes/graphql"
)

func main() {
	apiHandler := api.NewHandler()
	appHandler := app.NewHandler()
	authHandler := auth.NewHandler()
	gitHandler := git.NewHandler()
	graphqlHandler := graphql.NewHandler()

	subdomains := []handlers.SubdomainHandler{
		apiHandler,
		appHandler,
		authHandler,
		gitHandler,
		graphqlHandler,
	}

	defer func() { os.Exit(lib.Run(subdomains)) }()
}
