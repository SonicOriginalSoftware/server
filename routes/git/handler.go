package git

import (
	"api-server/lib/git"
	"api-server/lib/net/env"
	"api-server/lib/net/local"
	"io/ioutil"

	"fmt"
	"log"
	"net/http"
)

const prefix = "git"
const queryService = "service"

// Handler handles Git requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.outlog.Printf("[%v] %v %v\n", prefix, request.Method, request.URL)

	query := request.URL.Query()
	requestedService := query.Get(queryService)
	if requestedService != git.UploadService && requestedService != git.ReceiveService || requestedService == "" {
		http.Error(
			writer,
			fmt.Sprintf("Invalid service request: %v", query),
			http.StatusForbidden,
		)
		return
	}

	// TODO Do something with the body?
	_, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Cache-Control", "no-cache")
	writer.Header().Add("Content-Type", fmt.Sprintf("application/x-%v-advertisement", requestedService))
	// writer.Header().Add("Git-Protocol", "version=2")

	err = git.Execute(requestedService, request.URL.Path, writer)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
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

// NewHandler returns a new Handler
func NewHandler(outlog, errlog *log.Logger) *Handler {
	return &Handler{
		outlog: outlog,
		errlog: errlog,
	}
}
