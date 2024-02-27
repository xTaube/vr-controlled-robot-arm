package server

import (
	"fmt"
	"os"

	"github.com/xTaube/vr-controlled-robot-arm/video"
)

type CommandNotFound struct{}
func (err *CommandNotFound) Error() string {
	return "Command not found"
}

type CommandIdentifier byte

const (
	StartVideoStream CommandIdentifier = iota+1
	StopVideoStream
)

func startVideoStream() error {
	return video.StartFFMPEGVideoStreaming()
}

func stopVideoStream() error {
	return video.StopFFMPEGVideoStreaming()
}

func commandHandler(command CommandIdentifier) (string, error) {
	switch command {
	case StartVideoStream:
		err := startVideoStream()
		return fmt.Sprintf("Streaming to %s \n", os.Getenv("SERVER_ADDRESS")), err
	case StopVideoStream:
		err := stopVideoStream()
		return "Stream is down \n", err
	default:
		return "", &CommandNotFound{}
	}
}