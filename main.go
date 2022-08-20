//revive:disable:package-comments

package main

import (
	"os"
	"server/handlers"
	"server/lib"
	"server/routes/app"
	"server/routes/git"
	"server/routes/graphql"
)

func main() {
	appHandler := app.NewHandler()
	gitHandler := git.NewHandler()
	graphqlHandler := graphql.NewHandler()

	subdomains := []handlers.SubdomainHandler{
		appHandler,
		gitHandler,
		graphqlHandler,
	}

	defer func() { os.Exit(lib.Run(subdomains)) }()
}
