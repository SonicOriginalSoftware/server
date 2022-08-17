//revive:disable:package-comments

package git

import (
	"api-server/lib/net/env"
	"api-server/lib/net/local"
	"strings"

	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/transport"
	go_git "github.com/go-git/go-git/v5/plumbing/transport/server"
)

const prefix = "git"
const infoRefsService = "refs"
const receiveService = "receive-pack"
const uploadService = "upload-pack"

// Handler handles Git requests
type Handler struct {
	outlog *log.Logger
	errlog *log.Logger
	server transport.Transport
}

func (handler *Handler) handleError(writer http.ResponseWriter, errCode int, err error) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	handler.errlog.Printf("%s", err)
	http.Error(writer, err.Error(), errCode)
}

// ServeHTTP fulfills the http.Handler contract for Handler
func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Cache-Control", "no-cache")

	path := request.URL.Path
	err := fmt.Errorf("Invalid request: %v", path)

	pathParts := strings.Split(path, "/")
	service := pathParts[len(pathParts)-1]

	if service != receiveService && service != uploadService && service != infoRefsService {
		handler.handleError(writer, http.StatusForbidden, err)
		return
	}

	serviceQuery := request.URL.Query().Get("service")
	repoPath := strings.Join(pathParts[0:len(pathParts)-2], "/")
	endpoint := fmt.Sprintf("%v://%v%v", "https", request.Host, repoPath)
	transportEndpoint, err := transport.NewEndpoint(endpoint)
	if err != nil {
		handler.handleError(writer, http.StatusBadRequest, err)
		return
	}

	switch service {
	case infoRefsService:
		writer.Header().Set("Content-Type", fmt.Sprintf("application/x-%v-advertisement", serviceQuery))
		var session transport.UploadPackSession
		if session, err = handler.server.NewUploadPackSession(transportEndpoint, nil); err == nil {
			var refs *packp.AdvRefs
			if refs, err = session.AdvertisedReferencesContext(request.Context()); err == nil {
				refs.Prefix = [][]byte{[]byte("# service=git-upload-pack"), pktline.Flush}
				err = refs.Encode(writer)
			}
		}
	case receiveService:
		writer.Header().Set("Content-Type", fmt.Sprintf("application/x-%v-result", service))
		var session transport.ReceivePackSession
		if session, err = handler.server.NewReceivePackSession(transportEndpoint, nil); err == nil {
			receivePackRequest := packp.NewReferenceUpdateRequest()
			session.ReceivePack(request.Context(), receivePackRequest)
		}
	case uploadService:
		writer.Header().Set("Content-Type", fmt.Sprintf("application/x-%v-result", service))
		var session transport.UploadPackSession
		if session, err = handler.server.NewUploadPackSession(transportEndpoint, nil); err == nil {
			uploadPackRequest := packp.NewUploadPackRequest()
			session.UploadPack(request.Context(), uploadPackRequest)
		}
	}

	if err != nil {
		handler.handleError(writer, http.StatusBadRequest, err)
	}
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
		server: go_git.DefaultServer,
	}
}
