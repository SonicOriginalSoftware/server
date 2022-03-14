package git

import (
	"api-server/lib/git"
	"api-server/lib/net/env"
	"api-server/lib/net/local"
	"io"
	"os"
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

func handleInfoRefsRequest(service, repoPath string, writer http.ResponseWriter) {
	switch service {
	case git.ReceiveService, git.UploadService:
		writer.Header().Add("Content-Type", fmt.Sprintf("application/x-%v-advertisement", service))

		if err := git.InfoRefs(service, repoPath, writer); err != nil {
			http.Error(writer, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		}
	default:
		http.Error(
			writer,
			fmt.Sprintf("Invalid service request: %v", repoPath),
			http.StatusForbidden,
		)
	}
}

func handleServiceRequest(body io.ReadCloser, service, repoPath string, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", fmt.Sprintf("application/x-%v-result", service))

	if err := git.PackRequest(service, repoPath, body, io.MultiWriter(os.Stdout, writer)); err != nil {
		http.Error(writer, fmt.Sprintf("%s", err), http.StatusBadRequest)
	}
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handler.outlog.Printf("[%v] %v %v\n", prefix, request.Method, request.URL)

	writer.Header().Add("Cache-Control", "no-cache")

	path := request.URL.Path

	if strings.HasSuffix(path, git.InfoRefsPath) {
		handleInfoRefsRequest(request.URL.Query().Get(queryService), strings.TrimSuffix(path, git.InfoRefsPath), writer)
		return
	}

	pathParts := strings.Split(path, "/")
	service := pathParts[len(pathParts)-1]
	if service == git.ReceiveService || service == git.UploadService {
		handleServiceRequest(request.Body, service, strings.Join(pathParts[0:len(pathParts)-1], "/"), writer)
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
