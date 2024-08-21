package server

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/xTaube/vr-controlled-robot-arm/robot"
)

type WorkflowAbortedError struct {
	workflow_id string
	reason      string
}

func (e *WorkflowAbortedError) Error() string {
	return fmt.Sprintf("%s workflow was aborted. Reason: %s", e.workflow_id, e.reason)
}

type Workflow interface {
	Start()
}

type Step interface {
	Execute() error
	Revert() error
}

type PrepareRobotForCalibrationStep struct {
	workflow_id string
	robot       *robot.Robot
}

func (s *PrepareRobotForCalibrationStep) Execute() error {
	err := s.robot.StartCalibration()

	if err != nil {
		return &WorkflowAbortedError{s.workflow_id, err.Error()}
	}
	return nil
}

func (s *PrepareRobotForCalibrationStep) Revert() error {
	err := s.robot.AbortCalibration()

	if err != nil {
		return &WorkflowAbortedError{s.workflow_id, err.Error()}
	}
	return nil
}

type FinishRobotForCalibrationStep struct {
	workflow_id string
	robot       *robot.Robot
}

func (s *FinishRobotForCalibrationStep) Execute() error {
	err := s.robot.FinishCalibration()

	if err != nil {
		return &WorkflowAbortedError{s.workflow_id, err.Error()}
	}
	return nil
}

func (s *FinishRobotForCalibrationStep) Revert() error {
	return nil
}

type XYZAxisCalibrationStep struct {
	workflow_id string
	connection  *websocket.Conn
	robot       *robot.Robot
}

func (s *XYZAxisCalibrationStep) Execute() error {
	response := ResponseWithStringArguments{
		Code: RESPONSE_OK,
		Args: []string{"You're calibrating XYZ axis. Send '1' to confirm, send '2' to abort, send '3${X-deg}${Y-deg}${Z-deg}${V-deg}${W-deg}' to move."}}
	s.connection.WriteMessage(websocket.TextMessage, response.Parse())

	for {
		_, request, err := s.connection.ReadMessage()
		if err != nil {
			response := ErrorResponse{Code: RESPONSE_UNKNOWN_ERROR, Err: err}
			s.connection.WriteMessage(websocket.TextMessage, response.Parse())
		}
		command, args := ParseRequestArguments(string(request))

		switch command {
		case 1:
			if !s.robot.IsIdle() {
				response := ErrorResponse{
					Code: RESPONSE_ROBOT_CANNOT_EXECUTE_COMMAND_ERROR,
					Err:  &robot.RobotError{Code: robot.ROBOT_IS_IN_MOVE_ERROR, Err: nil},
				}
				s.connection.WriteMessage(websocket.TextMessage, response.Parse())
				continue
			}
			return nil

		case 2:
			return &WorkflowAbortedError{s.workflow_id, "user input"}

		case 3:
			fallback, err := s.robot.Move(
				robot.JointsAngles{Z: readFloat32(args[0]), Y: readFloat32(args[1]), X: readFloat32(args[2]), V: readFloat32(args[3]), W: readFloat32(args[4])},
			)
			if err != nil {
				response := ErrorResponse{Code: RESPONSE_ROBOT_CANNOT_EXECUTE_COMMAND_ERROR, Err: err}
				s.connection.WriteMessage(websocket.TextMessage, response.Parse())
				continue
			}
			response := ResponseWithFloat32Arguments{Code: RESPONSE_OK, Args: []float32{fallback.X, fallback.Y, fallback.Z}}
			s.connection.WriteMessage(websocket.TextMessage, response.Parse())

		default:
			response := ErrorResponse{
				Code: RESPONSE_UNKNOWN_COMMAND_ERROR,
				Err:  &CommandNotFound{command},
			}
			s.connection.WriteMessage(websocket.TextMessage, response.Parse())
		}
	}
}

func (s *XYZAxisCalibrationStep) Revert() error {
	return nil
}

type RobotCalibrationWorkflow struct {
	workflow_id string
	steps       []Step
	stepIdx     int
}

func (w *RobotCalibrationWorkflow) executeStep() error {
	err := w.steps[w.stepIdx].Execute()
	if err != nil {
		return err
	}
	return nil
}

func (w *RobotCalibrationWorkflow) revertStep() error {
	err := w.steps[w.stepIdx].Revert()
	if err != nil {
		return err
	}
	return nil
}

func (w *RobotCalibrationWorkflow) revert() {
	for {
		w.revertStep()
		w.stepIdx--

		if w.stepIdx < 0 {
			break
		}
	}
}

func (w *RobotCalibrationWorkflow) Start() error {
	w.stepIdx = 0
	for {
		err := w.executeStep()
		if err != nil {
			w.revert()
			return err
		}

		w.stepIdx++
		if w.stepIdx >= len(w.steps) {
			return nil
		}
	}
}

func InitRobotCalibrationWorkflow(connection *websocket.Conn, robot *robot.Robot) *RobotCalibrationWorkflow {
	workflow_id := "XYZ robot calibration"
	return &RobotCalibrationWorkflow{
		workflow_id: workflow_id,
		steps: []Step{
			&PrepareRobotForCalibrationStep{workflow_id: workflow_id, robot: robot},
			&XYZAxisCalibrationStep{workflow_id: workflow_id, connection: connection, robot: robot},
			&FinishRobotForCalibrationStep{workflow_id: workflow_id, robot: robot},
		},
	}
}
