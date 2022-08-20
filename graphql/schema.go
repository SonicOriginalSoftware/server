//revive:disable:package-comments

package graphql

import "github.com/graphql-go/graphql"

// Schema with rootQuery and rootMutation
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    &graphql.Object{},
	Mutation: &graphql.Object{},
})
