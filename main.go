package main

import (
	"fmt"
	"net/http"
	"os"
)

func root(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(os.Stdout, "Received a request!")
}

func main() {
	http.HandleFunc("/", root)

	address := ":8080"

	certPath := "./cert.pem"
	keyPath := "./key.pem"

	error := http.ListenAndServeTLS(address, certPath, keyPath, nil)
	if error != nil {
		fmt.Fprint(os.Stderr, error)
	}
}
