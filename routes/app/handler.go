package app

import (
	"api-server/lib/net/env"
	"api-server/lib/net/local"

	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const prefix = "app"

// Handler handles App requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger

	servePath string
}

//go:embed 404.html
var notFoundFile []byte

const defaultServePath = "public"
const indexFileName = "index.html"
const indexFileLength = len(indexFileName) - 1

func (handler *Handler) notFound(writer http.ResponseWriter, resource string, servePath string) {
	handler.errlog.Printf("[%v] Could not read resource at: %v\n", prefix, resource)

	indexStartIndex := len(resource) - 1 - indexFileLength
	if indexStartIndex > 0 && resource[indexStartIndex:] == indexFileName {
		writer.WriteHeader(http.StatusNotFound)
		if _, err := writer.Write(notFoundFile); err != nil {
			handler.errlog.Printf("[%v] %v", prefix, err)
			http.Error(writer, fmt.Sprintf("Could not retrieve %v", resource), http.StatusInternalServerError)
		}
		return
	}

	http.Error(writer, "Resource Not Found", http.StatusNotFound)
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	resourcePath := request.URL.Path
	if filepath.Ext(resourcePath) == "" {
		resourcePath = fmt.Sprintf("%v/%v", strings.TrimSuffix(resourcePath, "/"), indexFileName)
	}

	handler.outlog.Printf("[%v] request received: %v\n", prefix, request.URL)

	response, err := os.ReadFile(fmt.Sprintf("%v/%v", handler.servePath, resourcePath))
	if err != nil {
		handler.notFound(writer, resourcePath, handler.servePath)
		return
	}

	if _, err = writer.Write(response); err != nil {
		handler.errlog.Printf("[%v] Could not write response: %v", prefix, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
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
	servePath, isSet := os.LookupEnv("APP_SERVE_PATH")
	if !isSet {
		servePath = defaultServePath
	}

	return &Handler{
		outlog:    outlog,
		errlog:    errlog,
		servePath: servePath,
	}
}
