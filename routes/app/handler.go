package app

import (
	"fmt"
	"net/http"
	"os"
)

func notFound(writer http.ResponseWriter, request *http.Request) {
	notFoundMessage := "Resource Not Found"

	// FIXME Respond with a custom error page as well

	http.Error(writer, notFoundMessage, http.StatusNotFound)
}

func resource(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(os.Stdout, "Received an app root resource request: %v\n", request.URL.Path)

	servePath, isSet := os.LookupEnv("SERVE_PATH")
	if !isSet {
		servePath = "public"
	}

	resourcePath := request.URL.Path

	if request.URL.Path == "/" || request.URL.Path == "/app" || request.URL.Path == "/index.html" {
		resourcePath = "index.html"
	}

	response, error := os.ReadFile(fmt.Sprintf("%v/%v", servePath, resourcePath))
	if error != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not read resource at: %v%v", servePath, resourcePath))
		notFound(writer, request)
		return
	}

	writer.Write(response)
}

// Setup the request paths to app resources
func Setup() {
	http.HandleFunc("/", resource)
}
