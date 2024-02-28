package server

import (
	"fmt"

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

type CommandHandler struct {
	videoStream *video.VideoStream
}

func (ch *CommandHandler) Handle(command CommandIdentifier) (string, error) {
	switch command {
		case StartVideoStream:
			rtspServerAddress, err := ch.videoStream.Start()
			return fmt.Sprintf("Streaming to %s", rtspServerAddress), err
		
		case StopVideoStream:
			err := ch.videoStream.Stop()
			return "Stream is down", err
		
		default:
			return "", &CommandNotFound{}
	}
}

func InitCommandHandler(videoStream *video.VideoStream) *CommandHandler {
	return &CommandHandler{videoStream: videoStream}
}