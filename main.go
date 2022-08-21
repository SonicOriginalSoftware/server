//revive:disable:package-comments

package main

import (
	"server/handlers"
	"server/lib"
	"server/routes/app"
	"server/routes/git"
	"server/routes/graphql"
	"server/routes/grpc"

	"os"
)

func main() {
	appHandler := app.New()
	gitHandler := git.New()
	graphqlHandler := graphql.New()
	grpcHandler := grpc.New()

	subdomains := []handlers.SubdomainHandler{
		appHandler,
		gitHandler,
		graphqlHandler,
		grpcHandler,
	}

	defer func() { os.Exit(lib.Run(subdomains)) }()
}
