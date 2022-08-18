//revive:disable:package-comments

package api

import (
	"encoding/json"
	"server/lib/env"
	"server/lib/logging"
	"server/lib/net/local"
	"strings"

	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

const prefix = "api"

type postData struct {
	Variables map[string]interface{} `json:"variables"`
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
}

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    &graphql.Object{},
	Mutation: &graphql.Object{},
})

// Handler handles API requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

func (handler *Handler) handleGraphQL(writer http.ResponseWriter, request *http.Request) {
	var jsonData postData
	if err := json.NewDecoder(request.Body).Decode(&jsonData); err != nil {
		writer.WriteHeader(400)
		return
	}

	result := graphql.Do(graphql.Params{
		Context:        request.Context(),
		Schema:         schema,
		RequestString:  jsonData.Query,
		VariableValues: jsonData.Variables,
		OperationName:  jsonData.Operation,
	})

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		handler.errlog.Printf("Could not write result to response: %s", err)
	}
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.URL.Path, "/graphql") {
		handler.handleGraphQL(writer, request)
		return
	}

	http.Error(writer, "Not yet implemented!", http.StatusNotImplemented)
}

// Prefix is the subdomain prefix
func (handler *Handler) Prefix() string {
	return prefix
}

// Address returns the address the Handler will service
func (handler *Handler) Address() string {
	return env.Address(prefix, fmt.Sprintf("%v.%v", prefix, local.Path("")))
}

// NewHandler returns a new Handler
func NewHandler() *Handler {
	return &Handler{
		outlog: logging.NewLog(prefix),
		errlog: logging.NewError(prefix),
	}
}
