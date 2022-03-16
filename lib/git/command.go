package git

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	// InfoRefsPath is the path used when querying discovery info
	InfoRefsPath = "/info/refs"

	flushPacket         = "0000"
	gitCommand          = "git"
	uploadService       = "upload-pack"
	receiveService      = "receive-pack"
	statelessRPCOption  = "--stateless-rpc"
	advertiseRefsOption = "--advertise-refs"
	pktBufferSizeLength = 5
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

func trimRPC(fullServiceName string) string {
	return strings.TrimPrefix(fullServiceName, fmt.Sprintf("%v-", gitCommand))
}

func writePacketLine(writer io.Writer, service string) (bytes int, err error) {
	pktLine := fmt.Sprintf("# service=%v", service)
	pktLineSize := fmt.Sprintf("%04x", len(pktLine)+pktBufferSizeLength)

	if bytes, err = fmt.Fprintf(writer, "%v%v\n", pktLineSize, pktLine); bytes == 0 {
		err = fmt.Errorf("Could not write pkt-line to status buffer")
	}

	if bytes, err = fmt.Fprintf(writer, "%v", flushPacket); bytes == 0 {
		err = fmt.Errorf("Could not write flush packet to status buffer")
	}

	return
}

// InfoRefs returns the special info/refs file data from the requested repoPath
func InfoRefs(
	service, repoPath string,
	writer io.Writer,
) (cancel context.CancelFunc, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(
		ctx,
		gitCommand,
		trimRPC(service),
		statelessRPCOption,
		advertiseRefsOption,
		repoPath,
	)

	cmdPipe, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	if _, err = writePacketLine(writer, service); err != nil {
		return
	}

	if _, err = io.Copy(writer, cmdPipe); err != nil {
		return
	}

	return nil, cmd.Wait()
}

// PackRequest returns upload or receive pack info for a client request
func PackRequest(
	service,
	repoPath string,
	body io.Reader,
	writer io.Writer,
) (cancel context.CancelFunc, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(
		ctx,
		gitCommand,
		trimRPC(service),
		statelessRPCOption,
		repoPath,
	)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	if _, err = io.Copy(inPipe, body); err != nil {
		return
	}

	if _, err = io.Copy(writer, outPipe); err != nil {
		return
	}

	// FIXME On some occassions Wait() is returning an error with 0 len
	return nil, cmd.Wait()
}
