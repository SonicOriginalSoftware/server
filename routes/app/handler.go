package app

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler handles App requests
type Handler struct{}

//go:embed 404.html
var notFoundFile []byte

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

	fmt.Fprintf(os.Stdout, "Received an app root resource request: %v\n", resourcePath)

	servePath, isSet := os.LookupEnv("SERVE_PATH")
	if !isSet {
		servePath = "public"
	}

	response, error := os.ReadFile(fmt.Sprintf("%v/%v", servePath, resourcePath))
	if error != nil {
		notFound(writer, resourcePath, servePath)
		return
	}

	if _, error = writer.Write(response); error != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not write response: %v", error))
		http.Error(writer, error.Error(), http.StatusInternalServerError)
	}
}
