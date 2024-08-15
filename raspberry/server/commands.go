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
	MOVE_ROBOT
	SET_ROBOT_SPEED
	GET_ROBOT_CURRENT_POSITION
	CALIBRATE_ROBOT
)

type CommandHandler struct {
	videoStream              *video.VideoStream
	robot                    *robot.Robot
	robotCalibrationWorkflow *RobotCalibrationWorkflow
}

func (ch *CommandHandler) Handle(command_id CommandIdentifier, args []string) Response {
	log.Printf("Incoming command identitfier: %d\n", command_id)
	switch command_id {
	case START_VIDEO_STREAM:
		return ch.startVideoStreamCommandHandler()

	case STOP_VIDEO_STREAM:
		return ch.stopVideoStreamCommandHandler()

	case MOVE_ROBOT:
		return ch.moveArmCommandHandler(args)

	case SET_ROBOT_SPEED:
		return ch.setRobotSpeedCommandHandler(args)

	case GET_ROBOT_CURRENT_POSITION:
		return ch.getRobotCurrentPositionCommandHandler()

	case CALIBRATE_ROBOT:
		return ch.calibrateRobotCommandHandler()

	default:
		return &ErrorResponse{Code: RESPONSE_UNKNOWN_COMMAND_ERROR, Err: &CommandNotFound{command_id}}
	}
}

func (ch *CommandHandler) moveArmCommandHandler(command_args []string) Response {
	if len(command_args) < 5 {
		return &ErrorResponse{Code: RESPONSE_INVALID_PARAMETERS_NUMBER_ERROR}
	}
	log.Printf("Attempt to move robot by translation: [%s].\n", strings.Join(command_args, ", "))
	result, err := ch.robot.Move(
		robot.JointsAngles{
			Z: readFloat32(command_args[0]),
			Y: readFloat32(command_args[1]),
			X: readFloat32(command_args[2]),
			V: readFloat32(command_args[3]),
			W: readFloat32(command_args[4]),
		},
	)
	if err != nil {
		return &ErrorResponse{Code: RESPONSE_ROBOT_CANNOT_EXECUTE_COMMAND_ERROR, Err: err}
	}
	log.Println("Attempt finished.")
	return &ResponseWithFloat32Arguments{
		Code: RESPONSE_OK,
		Args: []float32{result.X, result.Y, result.Z, result.V, result.W},
	}
}

func (ch *CommandHandler) startVideoStreamCommandHandler() Response {
	log.Println("Turning stream on...")
	rtspServerAddress, err := ch.videoStream.Start()
	if err != nil {
		log.Printf("Error occured during turning stream on: %s\n", err)
		return &ErrorResponse{Code: RESPONSE_STREAM_ERROR, Err: err}
	}

	log.Printf("Streaming to %s.\n", rtspServerAddress)
	return &ResponseWithStringArguments{RESPONSE_OK, []string{rtspServerAddress}}
}

func (ch *CommandHandler) stopVideoStreamCommandHandler() Response {
	log.Println("Shutting off the stream...")
	err := ch.videoStream.Stop()
	if err != nil {
		log.Printf("Error occured during turning stream off: %s\n", err)
		return &ErrorResponse{Code: RESPONSE_STREAM_ERROR, Err: err}
	}
	log.Println("Stream stopped.")
	return &BaseResponse{Code: RESPONSE_OK}
}

func (ch *CommandHandler) setRobotSpeedCommandHandler(command_args []string) Response {
	if len(command_args) < 1 {
		return &ErrorResponse{Code: RESPONSE_INVALID_PARAMETERS_NUMBER_ERROR}
	}
	log.Printf("Attempt to set new robot speed: [%s].\n", strings.Join(command_args, ", "))
	err := ch.robot.SetSpeed(readFloat32(command_args[0]))
	if err != nil {
		return &ErrorResponse{Code: RESPONSE_ROBOT_CANNOT_EXECUTE_COMMAND_ERROR, Err: err}
	}

	log.Println("Attempt finished.")
	return &BaseResponse{Code: RESPONSE_OK}
}

func (ch *CommandHandler) getRobotCurrentPositionCommandHandler() Response {
	log.Println("Attempt to get current robot position.")
	currentPosition, err := ch.robot.GetCurrentPosition()
	if err != nil {
		return &ErrorResponse{Code: RESPONSE_ROBOT_CANNOT_EXECUTE_COMMAND_ERROR, Err: err}
	}
	log.Println("Attempt finished.")
	return &ResponseWithFloat32Arguments{
		Code: RESPONSE_OK,
		Args: []float32{currentPosition.X, currentPosition.Y, currentPosition.Z, currentPosition.V, currentPosition.W},
	}
}

func (ch *CommandHandler) calibrateRobotCommandHandler() Response {
	err := ch.robotCalibrationWorkflow.Start()
	if err != nil {
		return &ErrorResponse{Code: RESPONSE_ROBOT_CALIBRATION_ERROR, Err: err}
	}
	return &BaseResponse{Code: RESPONSE_OK}
}

func InitCommandHandler(
	videoStream *video.VideoStream,
	robot *robot.Robot,
	robotCalibrationWorkflow *RobotCalibrationWorkflow,
) *CommandHandler {
	return &CommandHandler{
		videoStream:              videoStream,
		robot:                    robot,
		robotCalibrationWorkflow: robotCalibrationWorkflow,
	}
}
