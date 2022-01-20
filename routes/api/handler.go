package api

import (
	"fmt"
	"net/http"
	"os"
)

func api(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(os.Stdout, "Received an API resource request!")
}

// Register the request paths to app resources
func Register() {
	http.HandleFunc("/api", api)
}
