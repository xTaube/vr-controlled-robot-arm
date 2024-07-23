package robot

import (
	"log"

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
	buff            []byte
	bytes_read      int
	terminationByte byte
}

func initUartBuffer() *UartBuffer {
	buffer := UartBuffer{
		buff:            make([]byte, UART_BUFFER_LEN),
		bytes_read:      0,
		terminationByte: byte(0x04),
	}

	return &buffer
}

func (ub *UartBuffer) load(port serial.Port) error {
	log.Println("UART: Trying to read...")
	totalBytes := 0
	for {
		tempBuff := make([]byte, 1)
		_, err := port.Read(tempBuff[:])
		if err != nil {
			return err
		}

		if tempBuff[0] == ub.terminationByte {
			break
		}

		ub.buff[totalBytes] = tempBuff[0]
		totalBytes++
	}
	log.Printf("UART: Read %d bytes.\n", totalBytes)
	ub.bytes_read = totalBytes
	return nil
}

func (ub *UartBuffer) Read() []byte {
	return ub.buff[:ub.bytes_read]
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
	log.Println("UART: Trying to send bytes...")
	data = append(data, u.buffer.terminationByte) // Add EOT byte
	if _, err := u.port.Write(data); err != nil {
		log.Printf("UART: writing data resulted in error: %s\n", err)
		return err
	}
	u.port.Drain()
	log.Printf("UART: Bytes sent %d\n", len(data))
	return nil
}

func (u *Uart) Get() ([]byte, error) {
	err := u.buffer.load(u.port)
	if err != nil {
		return nil, err
	}
	return u.buffer.Read(), nil
}
