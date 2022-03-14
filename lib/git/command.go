package git

import (
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

func cleanUpProcess(bytesWritten int64, cmd *exec.Cmd) (err error) {
	var errorMessage string
	if bytesWritten <= 0 || err != nil {
		if err := cmd.Process.Release(); err != nil {
			errorMessage = fmt.Sprintf("could not release process: %v", err.Error())
		}
		if err := cmd.Process.Kill(); err != nil {
			errorMessage = fmt.Sprintf("could not kill process: %v", err.Error())
		}
	}

	if bytesWritten <= 0 {
		defaultErrorMessage := "Could not write info/refs result"
		if errorMessage != "" {
			return fmt.Errorf("%v and %v", defaultErrorMessage, errorMessage)
		}
		return fmt.Errorf("%v", defaultErrorMessage)
	} else if err != nil {
		if errorMessage != "" {
			return fmt.Errorf("%v and %v", err.Error(), errorMessage)
		}
		return
	}

	return
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
func InfoRefs(service, repoPath string, writer io.Writer) (err error) {
	if _, err = writePacketLine(writer, service); err != nil {
		return
	}

	cmd := exec.Command(
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
		if exitErr, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf(string(exitErr.Stderr))
		}
		return
	}

	bytes, err := io.Copy(writer, cmdPipe)
	if err != nil || bytes <= 0 {
		return cleanUpProcess(bytes, cmd)
	}

	return cmd.Wait()
}

// PackRequest returns upload or receive pack info for a client request
func PackRequest(service, repoPath string, body io.Reader, writer io.Writer) (err error) {
	cmd := exec.Command(gitCommand, trimRPC(service), statelessRPCOption, repoPath)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf(string(exitErr.Stderr))
		}
		return
	}

	bytes, err := io.Copy(inPipe, body)
	if err != nil || bytes <= 0 {
		return cleanUpProcess(bytes, cmd)
	}

	bytes, err = io.Copy(writer, outPipe)
	if err != nil || bytes <= 0 {
		return cleanUpProcess(bytes, cmd)
	}

	return cmd.Wait()
}
