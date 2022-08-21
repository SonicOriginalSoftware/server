//revive:disable:package-comments

package graphql

import (
	"server/env"
	"server/graphql"
	"server/logging"
	"server/net/local"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	go_graphql "github.com/graphql-go/graphql"
)

const prefix = "graphql"

type postData struct {
	Variables map[string]interface{} `json:"variables"`
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
}

// Handler handles GraphQL API requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var jsonData postData
	if err := json.NewDecoder(request.Body).Decode(&jsonData); err != nil {
		writer.WriteHeader(400)
		return
	}

	result := go_graphql.Do(go_graphql.Params{
		Context:        request.Context(),
		Schema:         graphql.Schema,
		RequestString:  jsonData.Query,
		VariableValues: jsonData.Variables,
		OperationName:  jsonData.Operation,
	})

	if err := json.NewEncoder(writer).Encode(result); err != nil {
		handler.errlog.Printf("Could not write result to response: %s", err)
	}
}

// Prefix is the subdomain prefix
func (handler *Handler) Prefix() string {
	return prefix
}

// Address returns the address the Handler will service
func (handler *Handler) Address() string {
	return env.Address(prefix, fmt.Sprintf("%v.%v", prefix, local.Path("")))
}

// New returns a new Handler
func New() *Handler {
	return &Handler{
		outlog: logging.NewLog(prefix),
		errlog: logging.NewError(prefix),
	}
}
