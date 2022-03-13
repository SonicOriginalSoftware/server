package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	gitCommand     = "git"
	infoRefsSuffix = "/info/refs"
	uploadService  = "upload-pack"
	receiveService = "receive-pack"
	flushPacket    = "0000"
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

func writeFlushLine(writer *io.Writer) (bytes int, err error) {
	if bytes, err = fmt.Fprintf(*writer, "\n%v\n", flushPacket); bytes == 0 {
		err = fmt.Errorf("Could not write flush packet to status buffer")
	}

	return
}

func writePacketLine(writer *io.Writer, service string) (bytes int, err error) {
	pktLine := fmt.Sprintf("# service=%v", service)
	pktLine = fmt.Sprintf("%v%v\\n", fmt.Sprintf("%04x", len(pktLine)+5), pktLine)

	if bytes, err = (*writer).Write([]byte(pktLine)); bytes == 0 {
		err = fmt.Errorf("Could not write pkt-line to status buffer")
	}

	return
}

func writeServiceOutput(writer *io.Writer, output []byte) (bytes int, err error) {
	if bytes, err = fmt.Fprintf(*writer, "%s", output); bytes == 0 {
		err = fmt.Errorf("Could not write output to status buffer")
	}

	return
}

// Execute a git service command
//
// FIXME I believe the packet-line here is not performing as expected
// See git/connect.c line 339 and git/pkt-line.c line 408
func Execute(service string, repoPath string, statusWriter io.Writer) (err error) {
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
			err = fmt.Errorf(string(exitErr.Stderr))
		}
		return
	}

	writer := io.MultiWriter(statusWriter, os.Stdout)

	if _, err = writePacketLine(&writer, service); err != nil {
		return
	}

	if _, err = writeFlushLine(&writer); err != nil {
		return
	}

	_, err = writeServiceOutput(&writer, output)

	return
}
