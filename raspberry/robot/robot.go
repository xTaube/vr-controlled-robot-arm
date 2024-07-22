package robot

import (
	"encoding/binary"
	"log"
	"math"

	"go.bug.st/serial"
)

type RobotErrorCode uint8

const (
	ROBOT_INVALID_NUMBER_OF_PARAMETERS RobotErrorCode = iota + 10
	ROBOT_UNKNOWN_ACTION
	ROBOT_COMMUNICATION_ERROR
)

type RobotError struct {
	Code RobotErrorCode
	Err  error
}

func (err *RobotError) Error() string {
	switch err.Code {
	case ROBOT_INVALID_NUMBER_OF_PARAMETERS:
		return "Invalid number of parameters."
	case ROBOT_UNKNOWN_ACTION:
		return "Unknown action."
	default:
		return err.Err.Error()
	}
}

type ActionId uint8

const (
	MOVE_ACTION ActionId = iota + 1
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
)

type JointsTranslations struct {
	X float32
	Y float32
	Z float32
	V float32
	W float32
}

type Robot struct {
	uart *Uart
}

func (r *Robot) Move(jointsTranslations JointsTranslations) (*JointsTranslations, error) {
	data := make([]byte, W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE)

	data[ACTION_ID_OFFSET] = byte(MOVE_ACTION)
	binary.LittleEndian.PutUint32(
		data[X_JOINT_VALUE_OFFSET:X_JOINT_VALUE_OFFSET+X_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.X),
	)
	binary.LittleEndian.PutUint32(
		data[Y_JOINT_VALUE_OFFSET:Y_JOINT_VALUE_OFFSET+Y_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.Y),
	)
	binary.LittleEndian.PutUint32(
		data[Z_JOINT_VALUE_OFFSET:Z_JOINT_VALUE_OFFSET+Z_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.Z),
	)
	binary.LittleEndian.PutUint32(
		data[V_JOINT_VALUE_OFFSET:V_JOINT_VALUE_OFFSET+V_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.V),
	)
	binary.LittleEndian.PutUint32(
		data[W_JOINT_VALUE_OFFSET:W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.W),
	)

	err := r.uart.Send(data)
	if err != nil {
		return nil, &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	result, err := r.uart.Get()
	if err != nil {
		return nil, &RobotError{ROBOT_COMMUNICATION_ERROR, err}
	}

	result_code := result[0]
	log.Printf("Bytes received: %v\n", result)
	log.Println("UART received data: ")
	log.Printf("Code: %d\n", result_code)
	switch result_code {
	case byte(ROBOT_INVALID_NUMBER_OF_PARAMETERS):
		return nil, &RobotError{ROBOT_INVALID_NUMBER_OF_PARAMETERS, nil}
	case byte(ROBOT_UNKNOWN_ACTION):
		return nil, &RobotError{ROBOT_UNKNOWN_ACTION, nil}
	default:
		fallback := JointsTranslations{
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
