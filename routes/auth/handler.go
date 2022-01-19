package auth

import (
	"fmt"
	"net/http"
	"os"
)

func auth(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(os.Stdout, "Received an auth resource request!")
}

// Setup the request paths to app resources
func Setup() {
	http.HandleFunc("/auth", auth)
}
