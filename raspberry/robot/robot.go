package robot

import (
	"encoding/binary"
	"log"
	"math"

	"go.bug.st/serial"
)

type RobotCode uint8

const (
	ROBOT_OK = iota
	ROBOT_CALIBRATION_NEEDED
	ROBOT_COMMUNICATION_ERROR
)

type RobotActionResult struct {
	Code RobotCode
	Err  error
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

func (r *Robot) Move(jointsTranslations JointsTranslations) RobotActionResult {
	data := make([]byte, W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE)

	data[ACTION_ID_OFFSET] = byte(MOVE_ACTION)
	binary.BigEndian.PutUint32(
		data[X_JOINT_VALUE_OFFSET:X_JOINT_VALUE_OFFSET+X_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.X),
	)
	binary.BigEndian.PutUint32(
		data[Y_JOINT_VALUE_OFFSET:Y_JOINT_VALUE_OFFSET+Y_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.Y),
	)
	binary.BigEndian.PutUint32(
		data[Z_JOINT_VALUE_OFFSET:Z_JOINT_VALUE_OFFSET+Z_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.Z),
	)
	binary.BigEndian.PutUint32(
		data[V_JOINT_VALUE_OFFSET:V_JOINT_VALUE_OFFSET+V_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.V),
	)
	binary.BigEndian.PutUint32(
		data[W_JOINT_VALUE_OFFSET:W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE],
		math.Float32bits(jointsTranslations.W),
	)

	log.Printf("Bytes sent: \t %v \n", data)

	err := r.uart.Send(data)
	if err != nil {
		return RobotActionResult{Code: ROBOT_COMMUNICATION_ERROR, Err: err}
	}

	data1, err := r.uart.Get()
	if err != nil {
		return RobotActionResult{Code: ROBOT_COMMUNICATION_ERROR, Err: err}
	}

	log.Printf("Bytes received: \t %d \n", len(data1))

	log.Println("UART received data: ")
	log.Printf("Code: %d\n", data1[0])
	log.Printf("X: %f\n", math.Float32frombits(binary.BigEndian.Uint32(data1[X_JOINT_VALUE_OFFSET:X_JOINT_VALUE_OFFSET+X_JOINT_VALUE_SIZE])))
	log.Printf("Y: %f\n", math.Float32frombits(binary.BigEndian.Uint32(data1[Y_JOINT_VALUE_OFFSET:Y_JOINT_VALUE_OFFSET+Y_JOINT_VALUE_SIZE])))
	log.Printf("Z: %f\n", math.Float32frombits(binary.BigEndian.Uint32(data1[Z_JOINT_VALUE_OFFSET:Z_JOINT_VALUE_OFFSET+Z_JOINT_VALUE_SIZE])))
	log.Printf("V: %f\n", math.Float32frombits(binary.BigEndian.Uint32(data1[V_JOINT_VALUE_OFFSET:V_JOINT_VALUE_OFFSET+V_JOINT_VALUE_SIZE])))
	log.Printf("W: %f\n", math.Float32frombits(binary.BigEndian.Uint32(data1[W_JOINT_VALUE_OFFSET:W_JOINT_VALUE_OFFSET+W_JOINT_VALUE_SIZE])))

	return RobotActionResult{Code: RobotCode(data1[0]), Err: nil}
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
