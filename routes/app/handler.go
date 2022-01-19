package app

import (
	"fmt"
	"net/http"
	"os"
)

func root(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(os.Stdout, "Received an app root resource request!")
}

// Setup the request paths to app resources
func Setup() {
	http.HandleFunc("/", root)
}
