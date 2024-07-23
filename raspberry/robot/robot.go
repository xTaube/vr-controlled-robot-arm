package robot

import (
	"encoding/binary"
	"log"
	"math"

	"go.bug.st/serial"
)

type RobotErrorCode uint8

const (
	ROBOT_INVALID_NUMBER_OF_PARAMETERS_ERROR RobotErrorCode = iota + 10
	ROBOT_UNKNOWN_ACTION_ERROR
	ROBOT_NOT_CALIBRATED_ERROR
	ROBOT_SPEED_BEYOND_LIMIT_ERROR
	ROBOT_SPEED_TO_SLOW_ERROR
	ROBOT_IS_IN_MOVE_ERROR
	ROBOT_NOT_IN_CALIBRATION_MODE
	ROBOT_COMMUNICATION_ERROR
)

type RobotError struct {
	Code RobotErrorCode
	Err  error
}

func (err *RobotError) Error() string {
	switch err.Code {
	case ROBOT_INVALID_NUMBER_OF_PARAMETERS_ERROR:
		return "Invalid number of parameters."
	case ROBOT_UNKNOWN_ACTION_ERROR:
		return "Unknown action."
	case ROBOT_NOT_CALIBRATED_ERROR:
		return "Robot needs to be calibrated before operating."
	case ROBOT_SPEED_BEYOND_LIMIT_ERROR:
		return "Given speed is above possible max limit."
	case ROBOT_SPEED_TO_SLOW_ERROR:
		return "Given speed is below possible min limit."
	case ROBOT_IS_IN_MOVE_ERROR:
		return "Cannot perform action while robot is moving."
	case ROBOT_NOT_IN_CALIBRATION_MODE:
		return "Robot is not in calibration mode."
	default:
		if err.Err != nil {
			return err.Err.Error()
		}
		return ""
	}
}

type ActionId uint8

const (
	ACTION_MOVE ActionId = iota + 1
	ACTION_SET_SPEED
	ACTION_GET_CURRENT_POSITION
	ACTION_CHECK_ARM_CALIBRATION
	ACTION_START_CALIBARATION
	ACTION_FINISH_CALIBRATION
	ACTION_ABORT_CALIBRATION
	ACTION_CHECK_IDLE
)

const (
	ACTION_ID_OFFSET     uint8 = 0
	ACTION_ID_SIZE       uint8 = ACTION_ID_OFFSET + 1
	X_JOINT_VALUE_OFFSET uint8 = ACTION_ID_OFFSET + ACTION_ID_SIZE
	X_JOINT_VALUE_SIZE   uint8 = 4
	Y_JOINT_VALUE_OFFSET uint8 = X_JOINT_VALUE_OFFSET + X_JOINT_VALUE_SIZE
	Y_JOINT_VALUE_SIZE   uint8 = 4
	Z_JOINT_VALUE_OFFSET uint8 = Y_JOINT_VALUE_OFFSET + Y_JOINT_VALUE_SIZE
	Z_JOINT_VALUE_SIZE   uint8 = 4
	V_JOINT_VALUE_OFFSET uint8 = Z_JOINT_VALUE_OFFSET + Z_JOINT_VALUE_SIZE
	V_JOINT_VALUE_SIZE   uint8 = 4
	W_JOINT_VALUE_OFFSET uint8 = V_JOINT_VALUE_OFFSET + V_JOINT_VALUE_SIZE
	W_JOINT_VALUE_SIZE   uint8 = 4
	SPEED_VALUE_OFFSET   uint8 = ACTION_ID_OFFSET + ACTION_ID_SIZE
	SPEED_VALUE_SIZE     uint8 = 4
)

type JointsAngles struct {
	X float32
	Y float32
	Z float32
	V float32
	W float32
}

type Robot struct {
	uart *Uart
}

func (r *Robot) executeSimpleAction(action ActionId) error {
	data := make([]byte, ACTION_ID_OFFSET+ACTION_ID_SIZE)
	data[0] = byte(action)

	err := r.uart.Send(data)
	if err != nil {
		return &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	result, err := r.uart.Get()
	if err != nil {
		return &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	resultCode := RobotErrorCode(result[0])
	if resultCode >= 10 {
		return &RobotError{resultCode, nil}
	}
	return nil
}

func (r *Robot) Move(translations JointsAngles) (*JointsAngles, error) {
	data := make([]byte, W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE)

	data[ACTION_ID_OFFSET] = byte(ACTION_MOVE)
	binary.LittleEndian.PutUint32(
		data[X_JOINT_VALUE_OFFSET:X_JOINT_VALUE_OFFSET+X_JOINT_VALUE_SIZE],
		math.Float32bits(translations.X),
	)
	binary.LittleEndian.PutUint32(
		data[Y_JOINT_VALUE_OFFSET:Y_JOINT_VALUE_OFFSET+Y_JOINT_VALUE_SIZE],
		math.Float32bits(translations.Y),
	)
	binary.LittleEndian.PutUint32(
		data[Z_JOINT_VALUE_OFFSET:Z_JOINT_VALUE_OFFSET+Z_JOINT_VALUE_SIZE],
		math.Float32bits(translations.Z),
	)
	binary.LittleEndian.PutUint32(
		data[V_JOINT_VALUE_OFFSET:V_JOINT_VALUE_OFFSET+V_JOINT_VALUE_SIZE],
		math.Float32bits(translations.V),
	)
	binary.LittleEndian.PutUint32(
		data[W_JOINT_VALUE_OFFSET:W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE],
		math.Float32bits(translations.W),
	)

	err := r.uart.Send(data)
	if err != nil {
		return nil, &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	result, err := r.uart.Get()
	if err != nil {
		return nil, &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	resultCode := RobotErrorCode(result[0])
	log.Printf("Bytes received: %v\n", result)
	log.Println("UART received data: ")
	log.Printf("Code: %d\n", resultCode)

	if resultCode >= 10 {
		return nil, &RobotError{resultCode, nil}
	}

	fallback := JointsAngles{
		X: math.Float32frombits(binary.LittleEndian.Uint32(result[X_JOINT_VALUE_OFFSET : X_JOINT_VALUE_OFFSET+X_JOINT_VALUE_SIZE])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(result[Y_JOINT_VALUE_OFFSET : Y_JOINT_VALUE_OFFSET+Y_JOINT_VALUE_SIZE])),
		Z: math.Float32frombits(binary.LittleEndian.Uint32(result[Z_JOINT_VALUE_OFFSET : Z_JOINT_VALUE_OFFSET+Z_JOINT_VALUE_SIZE])),
		V: math.Float32frombits(binary.LittleEndian.Uint32(result[V_JOINT_VALUE_OFFSET : V_JOINT_VALUE_OFFSET+V_JOINT_VALUE_SIZE])),
		W: math.Float32frombits(binary.LittleEndian.Uint32(result[W_JOINT_VALUE_OFFSET : W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE])),
	}
	log.Printf("X: %f\n", fallback.X)
	log.Printf("Y: %f\n", fallback.Y)
	log.Printf("Z: %f\n", fallback.Z)
	log.Printf("V: %f\n", fallback.V)
	log.Printf("W: %f\n", fallback.W)
	return &fallback, nil
}

func (r *Robot) SetSpeed(speed float32) error {
	data := make([]byte, SPEED_VALUE_OFFSET+SPEED_VALUE_SIZE)
	data[0] = byte(ACTION_SET_SPEED)
	binary.LittleEndian.PutUint32(
		data[SPEED_VALUE_OFFSET:SPEED_VALUE_OFFSET+SPEED_VALUE_SIZE],
		math.Float32bits(speed),
	)

	err := r.uart.Send(data)
	if err != nil {
		return &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	result, err := r.uart.Get()
	if err != nil {
		return &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	resultCode := RobotErrorCode(result[0])
	if resultCode >= 10 {
		return &RobotError{resultCode, nil}
	}

	return nil
}

func (r *Robot) GetCurrentPosition() (*JointsAngles, error) {
	data := make([]byte, ACTION_ID_OFFSET+ACTION_ID_SIZE)
	data[0] = byte(ACTION_GET_CURRENT_POSITION)

	err := r.uart.Send(data)
	if err != nil {
		return nil, &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	result, err := r.uart.Get()
	if err != nil {
		return nil, &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}
	log.Printf("Bytes received: %v\n", result)

	resultCode := RobotErrorCode(result[0])
	if resultCode >= 10 {
		return nil, &RobotError{resultCode, nil}
	}

	currentPosition := JointsAngles{
		X: math.Float32frombits(binary.LittleEndian.Uint32(result[X_JOINT_VALUE_OFFSET : X_JOINT_VALUE_OFFSET+X_JOINT_VALUE_SIZE])),
		Y: math.Float32frombits(binary.LittleEndian.Uint32(result[Y_JOINT_VALUE_OFFSET : Y_JOINT_VALUE_OFFSET+Y_JOINT_VALUE_SIZE])),
		Z: math.Float32frombits(binary.LittleEndian.Uint32(result[Z_JOINT_VALUE_OFFSET : Z_JOINT_VALUE_OFFSET+Z_JOINT_VALUE_SIZE])),
		V: math.Float32frombits(binary.LittleEndian.Uint32(result[V_JOINT_VALUE_OFFSET : V_JOINT_VALUE_OFFSET+V_JOINT_VALUE_SIZE])),
		W: math.Float32frombits(binary.LittleEndian.Uint32(result[W_JOINT_VALUE_OFFSET : W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE])),
	}
	log.Println("Current position:")
	log.Printf("X: %f\n", currentPosition.X)
	log.Printf("Y: %f\n", currentPosition.Y)
	log.Printf("Z: %f\n", currentPosition.Z)
	log.Printf("V: %f\n", currentPosition.V)
	log.Printf("W: %f\n", currentPosition.W)

	return &currentPosition, nil
}

func (r *Robot) StartCalibration() error {
	return r.executeSimpleAction(ACTION_START_CALIBARATION)
}

func (r *Robot) FinishCalibration() error {
	return r.executeSimpleAction(ACTION_FINISH_CALIBRATION)
}

func (r *Robot) AbortCalibration() error {
	return r.executeSimpleAction(ACTION_ABORT_CALIBRATION)
}

func (r *Robot) IsIdle() bool {
	err := r.executeSimpleAction(ACTION_CHECK_IDLE)
	return err == nil
}

func (r *Robot) ShutDown() {
	r.uart.Close()
}

func InitRobot(uartConfig UartConfig) (*Robot, error) {
	log.Println("Initializing UART...")
	uart, err := initUart(
		uartConfig.PortName,
		&serial.Mode{
			BaudRate: uartConfig.BaudRate,
			DataBits: uartConfig.DataBits,
			Parity:   uartConfig.Parity,
			StopBits: uartConfig.StopBits,
		},
	)

	if err != nil {
		return nil, err
	}
	log.Println("UART initialized.")

	return &Robot{uart: uart}, nil
}
