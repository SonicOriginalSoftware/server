package auth

import (
	"fmt"
	"net/http"
	"os"
)

func auth(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(os.Stdout, "Received an auth resource request!")
	http.Error(writer, "Not yet implemented!", http.StatusNotImplemented)
}

// Register the request paths to app resources
func Register() {
	http.HandleFunc("/auth", auth)
}
