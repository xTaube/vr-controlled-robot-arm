package robot

import (
	"log"
	"slices"

	"go.bug.st/serial"
)

const UART_BUFFER_LEN = 256

type UartConfig struct {
	PortName string
	Parity   serial.Parity
	StopBits serial.StopBits
	DataBits int
	BaudRate int
}

type UartBuffer struct {
	buff      []byte
	bytesRead int
}

func initUartBuffer() *UartBuffer {
	buffer := UartBuffer{
		buff:      make([]byte, UART_BUFFER_LEN),
		bytesRead: 0,
	}

	return &buffer
}

func (ub *UartBuffer) load(port serial.Port) error {
	// log.Println("UART: Trying to read...")
	bytesToRead := make([]byte, 1)
	_, err := port.Read(bytesToRead[:])
	if err != nil {
		return err
	}

	totalBytes := 0
	for totalBytes < int(bytesToRead[0]) {
		n, err := port.Read(ub.buff[totalBytes:])
		if err != nil {
			return err
		}
		totalBytes += n
	}
	// log.Printf("UART: Read %d bytes.\n", totalBytes)
	ub.bytesRead = totalBytes
	return nil
}

func (ub *UartBuffer) Read() []byte {
	return ub.buff[:ub.bytesRead]
}

type Uart struct {
	portName string
	port     serial.Port
	buffer   *UartBuffer
}

func initUart(
	portName string,
	mode *serial.Mode,
) (*Uart, error) {
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

func (u *Uart) Send(data []byte) error {
	// log.Println("UART: Trying to send bytes...")
	data = slices.Insert(data, 0, byte(len(data))) // Add number of bytes to read at the beginning

	_, err := u.port.Write(data)
	if err != nil {
		log.Printf("UART: writing data resulted in error: %s\n", err)
		return err
	}

	u.port.Drain()
	// log.Printf("UART: Bytes sent %d\n", n)
	return nil
}

func (u *Uart) Get() ([]byte, error) {
	err := u.buffer.load(u.port)
	if err != nil {
		return nil, err
	}
	return u.buffer.Read(), nil
}
