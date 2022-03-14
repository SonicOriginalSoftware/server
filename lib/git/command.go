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

func writeServiceOutputLine(writer io.Writer, args []string) (bytes int64, err error) {
	cmd := exec.Command(gitCommand, args...)
	cmdPipe, err := cmd.StdoutPipe()

	if err = cmd.Start(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf(string(exitErr.Stderr))
		}
		return
	}

	if bytes, err = io.Copy(writer, cmdPipe); bytes <= 0 {
		err = fmt.Errorf("Could not write service output")
		return
	}

	err = cmd.Wait()

	return
}

// Execute a git service command
//
// FIXME I believe the packet-line here is not performing as expected
// Rough stack trace:
//   cmd_main (remote—curl.c)
//   parse_fetch
//   fetch
//   discover_refs
// 	 parse_git_refs
//     discover_version (connect.c)
//       packet_reader_peek (pkt-line.c)
//       packet_reader_read
//       packet_read_with_status
// Multiple points in this that return PACKET_READ_EOF.
// I believe we are getting through the initial case.
// There’s another case though that I think I was able to get to
// in certain buffer write orders, but I’m not sure why it was failing there
func Execute(service, repoPath string, advertiseRefs bool, writer io.Writer) (err error) {
	trimmedService := strings.TrimPrefix(service, fmt.Sprintf("%v-", gitCommand))
	repoPath = strings.TrimSuffix(repoPath, InfoRefsPath)

	args := []string{trimmedService}
	if trimmedService == uploadService {
		args = append(args, statelessRPCOption)
	}
	if advertiseRefs {
		args = append(args, advertiseRefsOption)
	}

	if _, err = writePacketLine(writer, service); err != nil {
		return
	}

	_, err = writeServiceOutputLine(writer, append(args, repoPath))

	return
}
