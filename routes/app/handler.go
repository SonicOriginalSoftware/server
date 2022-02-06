package app

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

func notFound(writer http.ResponseWriter, resource string, servePath string) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not read resource at: %v", resource))

	if indexStartIndex := len(resource) - 1 - indexFileLength; indexStartIndex > 0 && resource[indexStartIndex:] == indexFileName {
		writer.WriteHeader(http.StatusNotFound)
		if _, error := writer.Write(notFoundFile); error != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%v", error))
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

	handler.outlog.Printf("Received an app resource request: %v\n", resourcePath)

	response, error := os.ReadFile(fmt.Sprintf("%v/%v", handler.servePath, resourcePath))
	if error != nil {
		notFound(writer, resourcePath, handler.servePath)
		return
	}

	if _, error = writer.Write(response); error != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not write response: %v", error))
		http.Error(writer, error.Error(), http.StatusInternalServerError)
	}
}

// NewHandler returns a new Handler
func NewHandler(outlog, errlog *log.Logger) *Handler {
	servePath, isSet := os.LookupEnv("SERVE_PATH")
	if !isSet {
		servePath = defaultServePath
	}

	return &Handler{
		outlog:    outlog,
		errlog:    errlog,
		servePath: servePath,
	}
}
