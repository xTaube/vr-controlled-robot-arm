package robot

import (
	"log"

	"go.bug.st/serial"
)

const UART_BUFFER_LEN = 256

type NoDataAvailableError struct{}

func (err *NoDataAvailableError) Error() string {
	return "No data available on port"
}

type UartBuffer struct {
	buff       []byte
	bytes_read int
}

func initUartBuffer() *UartBuffer {
	buffer := UartBuffer{
		buff:       make([]byte, UART_BUFFER_LEN),
		bytes_read: 0,
	}

	return &buffer
}

func (ub *UartBuffer) load(port serial.Port) error {
	bytes_read, err := port.Read(ub.buff)

	if err != nil {
		return err
	}

	ub.bytes_read = bytes_read
	if bytes_read == 0 {
		return &NoDataAvailableError{}
	}

	return nil
}

func (ub *UartBuffer) get() []byte {
	return ub.buff[:ub.bytes_read]
}

type Uart struct {
	portName string
	port     serial.Port
	buffer   *UartBuffer
}

func InitUart(
	portName string,
	baudRate int,
	parity serial.Parity,
	dataBits int,
	stopBits serial.StopBits,
) (*Uart, error) {
	mode := &serial.Mode{
		BaudRate: baudRate,
		Parity:   parity,
		DataBits: dataBits,
		StopBits: stopBits,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}

	buffer := initUartBuffer()
	uart := Uart{
		portName: portName,
		port:     port,
		buffer:   buffer,
	}
	return &uart, nil
}

func (u *Uart) Close() error {
	err := u.port.Close()
	if err != nil {
		log.Printf("Failed to close UART port %s\n", err)
		return err
	}
	return nil
}

func (u *Uart) Send(message string) (string, error) {
	if _, err := u.port.Write([]byte(message)); err != nil {
		return "", err
	}

	if err := u.buffer.load(u.port); err != nil {
		return "", err
	}

	return string(u.buffer.get()), nil
}
