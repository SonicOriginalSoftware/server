package git

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	gitCommand     = "git"
	infoRefsSuffix = "/info/refs"
	uploadService  = "upload-pack"
	receiveService = "receive-pack"
)

var (
	// UploadService is the name of the git-upload-pack service
	UploadService string
	// ReceiveService is the name of the git-receive-pack service
	ReceiveService string

	uploadPackOptions []string
)

func init() {
	UploadService = fmt.Sprintf("%v-%v", gitCommand, uploadService)
	ReceiveService = fmt.Sprintf("%v-%v", gitCommand, receiveService)
	uploadPackOptions = []string{"--stateless-rpc", "--http-backend-info-refs"}
}

// Execute a git service command
func Execute(service string, repoPath string) (status string, err error) {
	trimmedService := strings.TrimPrefix(service, fmt.Sprintf("%v-", gitCommand))
	repoPath = strings.TrimSuffix(repoPath, infoRefsSuffix)

	args := []string{trimmedService}
	if trimmedService == uploadService {
		args = append(args, uploadPackOptions...)
	}
	args = append(args, repoPath)

	output, err := exec.Command(gitCommand, args...).Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status = string(exitErr.Stderr)
		} else {
			status = err.Error()
		}
		return
	}

	pktLineTrailer := fmt.Sprintf("# service=%v", service)
	status = fmt.Sprintf(
		"%v\n%v",
		fmt.Sprintf("%v%v", fmt.Sprintf("%04x", len(pktLineTrailer)+5), pktLineTrailer),
		string(output),
	)

	return
}
