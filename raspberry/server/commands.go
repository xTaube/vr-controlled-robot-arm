package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/xTaube/vr-controlled-robot-arm/robot"
	"github.com/xTaube/vr-controlled-robot-arm/video"
)

type CommandIdentifier byte

type CommandNotFound struct {
	code CommandIdentifier
}

func (err *CommandNotFound) Error() string {
	return fmt.Sprintf("Command with identifier: %d not found.", err.code)
}

const (
	START_VIDEO_STREAM CommandIdentifier = iota + 1
	STOP_VIDEO_STREAM
	MOVE_ARM
)

type CommandHandler struct {
	videoStream *video.VideoStream
	robot       *robot.Robot
}

func (ch *CommandHandler) Handle(command_id CommandIdentifier, args []string) ([]byte, error) {
	log.Printf("Incoming command identitfier: %d\n", command_id)
	switch command_id {
	case START_VIDEO_STREAM:
		return ch.startVideoStreamCommandHandler()

	case STOP_VIDEO_STREAM:
		return ch.stopVideoStreamCommandHandler()

	case MOVE_ARM:
		return ch.moveArmCommandHandler(args)

	default:
		return nil, &CommandNotFound{command_id}
	}
}

func (ch *CommandHandler) moveArmCommandHandler(command_args []string) ([]byte, error) {
	log.Printf("Attempt to move arm by translation: [%s].\n", strings.Join(command_args, ", "))
	result := ch.robot.Move(
		robot.JointsTranslations{
			X: readFloat32(command_args[0]),
			Y: readFloat32(command_args[1]),
			Z: readFloat32(command_args[2]),
			V: readFloat32(command_args[3]),
			W: readFloat32(command_args[4]),
		},
	)
	log.Println("Attempt finished.")
	return SuccessResponse(OK, ""), result.Err
}

func (ch *CommandHandler) startVideoStreamCommandHandler() ([]byte, error) {
	log.Println("Turning stream on...")
	rtspServerAddress, err := ch.videoStream.Start()
	if err != nil {
		return nil, err
	}

	log.Printf("Streaming to %s.\n", rtspServerAddress)
	return SuccessResponse(OK, fmt.Sprintf("Stream available on %s.", rtspServerAddress)), nil
}

func (ch *CommandHandler) stopVideoStreamCommandHandler() ([]byte, error) {
	log.Println("Shutting off the stream...")
	err := ch.videoStream.Stop()
	log.Println("Stream stopped.")
	return SuccessResponse(OK, "Stream disabled."), err
}

func InitCommandHandler(videoStream *video.VideoStream, robot *robot.Robot) *CommandHandler {
	return &CommandHandler{videoStream: videoStream, robot: robot}
}
