package server

import (
	"fmt"
	"log"

	"github.com/xTaube/vr-controlled-robot-arm/video"
)

type CommandNotFound struct{}

func (err *CommandNotFound) Error() string {
	return "Command not found"
}

type CommandIdentifier byte

const (
	StartVideoStream CommandIdentifier = iota + 1
	StopVideoStream
)

type CommandHandler struct {
	videoStream *video.VideoStream
}

func (ch *CommandHandler) Handle(command CommandIdentifier) (string, error) {
	log.Printf("Incoming: %d\n", command)
	switch command {
	case StartVideoStream:
		log.Println("Turning stream on...")
		rtspServerAddress, err := ch.videoStream.Start()
		if err != nil {
			return "", err
		}
		log.Printf("Streaming to %s", rtspServerAddress)
		return fmt.Sprintf("Stream available on %s", rtspServerAddress), nil

	case StopVideoStream:
		log.Println("Shutting off stream...")
		err := ch.videoStream.Stop()
		log.Println("Stream stopped")
		return "Stream disabled", err

	default:
		return "", &CommandNotFound{}
	}
}

func InitCommandHandler(videoStream *video.VideoStream) *CommandHandler {
	return &CommandHandler{videoStream: videoStream}
}
