package git

import (
	"api-server/lib/git"
	"api-server/lib/net/env"
	"api-server/lib/net/local"
	"io"
	"os"
	"os/exec"
	"strings"

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

func handleError(writer http.ResponseWriter, errCode int, err error) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(errCode)

	if exitErr, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf(string(exitErr.Stderr))
	}

	bytes, err := fmt.Fprintf(io.MultiWriter(writer, os.Stderr), "%s", err)
	if bytes <= 0 || err != nil {
		if bytes <= 0 {
			err = fmt.Errorf("Could not write error to error writer")
		}
		log.Fatalf("%v", err)
	}
}

func handleInfoRefsRequest(service, repoPath string, writer http.ResponseWriter) {
	if service != git.UploadService && service != git.ReceiveService {
		handleError(
			writer,
			http.StatusForbidden,
			fmt.Errorf("Invalid service request: %v", repoPath),
		)
		return
	}

	writer.Header().Set("Content-Type", fmt.Sprintf("application/x-%v-advertisement", service))

	if cancel, err := git.InfoRefs(service, repoPath, writer); err != nil {
		defer cancel()
		handleError(writer, http.StatusInternalServerError, err)
	}
}

func handleServiceRequest(body io.ReadCloser, service, repoPath string, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", fmt.Sprintf("application/x-%v-result", service))

	if cancel, err := git.PackRequest(service, repoPath, body, writer); err != nil {
		defer cancel()
		handleError(writer, http.StatusBadRequest, err)
	}
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.outlog.Printf("[%v] %v %v\n", prefix, request.Method, request.URL)

	writer.Header().Add("Cache-Control", "no-cache")

	path := request.URL.Path

	if strings.HasSuffix(path, git.InfoRefsPath) {
		handleInfoRefsRequest(
			request.URL.Query().Get(queryService),
			strings.TrimSuffix(path, git.InfoRefsPath),
			writer,
		)
		return
	}

	pathParts := strings.Split(path, "/")
	service := pathParts[len(pathParts)-1]
	if service == git.ReceiveService || service == git.UploadService {
		handleServiceRequest(
			request.Body,
			service,
			strings.Join(pathParts[0:len(pathParts)-1], "/"),
			writer,
		)
		return
	}

	http.Error(
		writer,
		fmt.Sprintf("Invalid request: %v", path),
		http.StatusForbidden,
	)
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
