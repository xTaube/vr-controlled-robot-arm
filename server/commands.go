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

func startVideoStream(
	devicePath string,
	resolution [2]int, 
	framerate int, 
	serverAddress string,
) error {
	builder := video.InitFFMPEGCommandBuilder()
	builder.SetStreamFormat(
		video.V4L2,
	).SetInputFormat(
		video.MJPEG,
	).SetVideoSize(
		resolution[0],
		resolution[1],
	).SetFramerate(
		framerate,
	).SetDevice(
		devicePath,
	).SetRTSPOutput(serverAddress)
	return builder.Execute()
}

func stopVideoStream() error {
	return video.StopFFMPEGVideoStreaming()
}

func commandHandler(command CommandIdentifier) (string, error) {
	switch command {
	case StartVideoStream:
		err := startVideoStream(
			os.Getenv("CAMERA_DEVICE_PATH"),
			[2]int{1920, 1080},
			30,
			os.Getenv("SERVER_ADDRESS"),
		)
		return fmt.Sprintf("Streaming to %s \n", os.Getenv("SERVER_ADDRESS")), err
	case StopVideoStream:
		err := stopVideoStream()
		return "Stream is down \n", err
	default:
		return "", &CommandNotFound{}
	}
}