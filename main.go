//revive:disable:package-comments

package main

import (
	"os"
	"server/handlers"
	"server/lib"
	"server/routes/app"
	"server/routes/git"
	"server/routes/graphql"
	"server/routes/grpc"
)

func main() {
	appHandler := app.NewHandler()
	gitHandler := git.NewHandler()
	graphqlHandler := graphql.NewHandler()
	grpcHandler := grpc.NewHandler()

	subdomains := []handlers.SubdomainHandler{
		appHandler,
		gitHandler,
		graphqlHandler,
		grpcHandler,
	}

	defer func() { os.Exit(lib.Run(subdomains)) }()
}
